package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb_benchmark "github.com/fffeng99999/hcp-server/api/generated/benchmark"
	pb_metric "github.com/fffeng99999/hcp-server/api/generated/metric"
	pb_node "github.com/fffeng99999/hcp-server/api/generated/node"
	pb_transaction "github.com/fffeng99999/hcp-server/api/generated/transaction"
	"github.com/fffeng99999/hcp-server/internal/config"
	"github.com/fffeng99999/hcp-server/internal/database"
	"github.com/fffeng99999/hcp-server/internal/grpc/handlers"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
	"github.com/fffeng99999/hcp-server/internal/service"
	"github.com/fffeng99999/hcp-server/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config", "configs", "Path to config directory")
	flag.Parse()

	// 1. Load Config
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Init Logger
	if err := utils.InitLogger(cfg.Log.Level); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer utils.Logger.Sync()
	utils.Logger.Info("Starting hcp-server", zap.String("version", "1.0.0"))

	// 3. Connect DB
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 3.1 Connect Redis
	rdb, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		utils.Logger.Warn("Failed to connect to redis, cache will be disabled", zap.Error(err))
		// Don't exit, just warn for now as it might be optional or dev env
	} else {
		utils.Logger.Info("Connected to Redis", zap.String("addr", cfg.Redis.Addr))
		defer rdb.Close()
	}

	// 4. Auto Migrate
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		utils.Logger.Info("Running migrations...")
		err = db.AutoMigrate(
			&models.Benchmark{},
			&models.Transaction{},
			&models.Node{},
			&models.Metric{},
			&models.Anomaly{},
		)
		if err != nil {
			utils.Logger.Fatal("Migration failed", zap.Error(err))
		}
		utils.Logger.Info("Migration completed successfully")
		return
	}

	// 5. Init Repositories
	benchmarkRepo := repository.NewBenchmarkRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	nodeRepo := repository.NewNodeRepository(db)
	metricRepo := repository.NewMetricRepository(db)

	// 6. Init Services
	benchmarkService := service.NewBenchmarkService(benchmarkRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	nodeService := service.NewNodeService(nodeRepo)
	metricService := service.NewMetricService(metricRepo)

	// 7. Init gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		utils.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()

	// Register Handlers
	benchmarkHandler := handlers.NewBenchmarkHandler(benchmarkService)
	pb_benchmark.RegisterBenchmarkServiceServer(s, benchmarkHandler)

	transactionHandler := handlers.NewTransactionHandler(transactionService)
	pb_transaction.RegisterTransactionServiceServer(s, transactionHandler)

	nodeHandler := handlers.NewNodeHandler(nodeService)
	pb_node.RegisterNodeServiceServer(s, nodeHandler)

	metricHandler := handlers.NewMetricHandler(metricService)
	pb_metric.RegisterMetricServiceServer(s, metricHandler)

	// 8. Start Server
	utils.Logger.Info("Server listening", zap.Int("port", cfg.Server.Port))

	go func() {
		if err := s.Serve(lis); err != nil {
			utils.Logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Logger.Info("Shutting down server...")
	s.GracefulStop()
	utils.Logger.Info("Server stopped")
}
