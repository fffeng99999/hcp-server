package repository

import (
	"context"
	"time"

	"github.com/fffeng99999/hcp-server/internal/models"
)

type BenchmarkRepository interface {
	Create(ctx context.Context, benchmark *models.Benchmark) error
	GetByID(ctx context.Context, id string) (*models.Benchmark, error)
	List(ctx context.Context, page, pageSize int) ([]models.Benchmark, int64, error)
	Update(ctx context.Context, benchmark *models.Benchmark) error
	Delete(ctx context.Context, id string) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	GetByHash(ctx context.Context, hash string) (*models.Transaction, error)
	List(ctx context.Context, filter TransactionFilter, page, pageSize int) ([]models.Transaction, int64, error)
	GetStats(ctx context.Context, benchmarkID string) (*TransactionStats, error)
}

type TransactionFilter struct {
	BenchmarkID string
	FromAddress string
	ToAddress   string
	Status      string
}

type TransactionStats struct {
	TotalTransactions int64
	PendingCount      int64
	ConfirmedCount    int64
	FailedCount       int64
	AvgLatencyMs      float64
}

type NodeRepository interface {
	Create(ctx context.Context, node *models.Node) error
	GetByID(ctx context.Context, id string) (*models.Node, error)
	List(ctx context.Context, filter NodeFilter, page, pageSize int) ([]models.Node, int64, error)
	Update(ctx context.Context, node *models.Node) error
	UpdateStatus(ctx context.Context, id, status string) error
}

type NodeFilter struct {
	Role   string
	Status string
	Region string
}

type MetricRepository interface {
	Create(ctx context.Context, metric *models.Metric) error
	CreateBatch(ctx context.Context, metrics []*models.Metric) error
	GetNodeMetrics(ctx context.Context, nodeID, metricName string, startTime, endTime time.Time, page, pageSize int) ([]models.Metric, int64, error)
	GetBenchmarkMetrics(ctx context.Context, benchmarkID, metricName string, page, pageSize int) ([]models.Metric, int64, error)
}
