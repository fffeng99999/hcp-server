package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

const (
	// Credentials from core memory
	dsn = "host=192.168.58.102 user=user_rbc3B8 password=password_DfA4Pw dbname=hcp_server port=5432 sslmode=disable"
)

func main() {
	log.Println("Connecting to database...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected successfully.")

	// 1. Run Migration
	migrationFile := "../../internal/database/migrations/001_init_schema.sql"
	absPath, _ := filepath.Abs(migrationFile)
	log.Printf("Reading migration file: %s", absPath)

	content, err := os.ReadFile(migrationFile)
	if err != nil {
		// Try absolute path based on known structure if relative fails
		migrationFile = "/home/hcp-dev/hcp-project/hcp-server/internal/database/migrations/001_init_schema.sql"
		content, err = os.ReadFile(migrationFile)
		if err != nil {
			log.Fatalf("Failed to read migration file: %v", err)
		}
	}

	// Split by semicolon might be too simple if functions have semicolons, but the schema file
	// uses $$ for functions, so internal semicolons are safe if we execute the whole block?
	// Postgres driver might not support multiple statements in one Exec.
	// Let's try executing the whole string first.
	// pq driver supports multiple statements.

	log.Println("Executing migration SQL...")
	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}
	log.Println("Migration completed.")

	// 2. Seed Data
	log.Println("Seeding data...")
	seedData(db)
	log.Println("Seeding completed.")
}

func seedData(db *sql.DB) {
	// 1. Nodes
	nodes := []struct {
		ID, Name, Address, Role, Status string
		TrustScore                      float64
	}{
		{"node-001", "Validator-One", "192.168.1.101:8080", "validator", "online", 98.5},
		{"node-002", "Validator-Two", "192.168.1.102:8080", "validator", "online", 99.1},
		{"node-003", "Observer-One", "192.168.1.103:8080", "observer", "online", 100.0},
	}

	for _, n := range nodes {
		_, err := db.Exec(`
			INSERT INTO nodes (id, name, address, role, status, trust_score, last_heartbeat)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT (id) DO NOTHING`,
			n.ID, n.Name, n.Address, n.Role, n.Status, n.TrustScore)
		if err != nil {
			log.Printf("Error seeding node %s: %v", n.ID, err)
		}
	}

	// 2. Benchmarks
	benchID := "550e8400-e29b-41d4-a716-446655440000"
	_, err := db.Exec(`
		INSERT INTO benchmarks (id, name, algorithm, node_count, duration, actual_tps, status, created_at)
		VALUES ($1, 'tPBFT-Baseline', 'tPBFT', 4, 60, 1500.50, 'completed', NOW() - INTERVAL '1 hour')
		ON CONFLICT (id) DO NOTHING`, benchID)
	if err != nil {
		log.Printf("Error seeding benchmark: %v", err)
	}

	// 3. Transactions
	// Manually create partition for current month to avoid trigger deadlock during insert
	log.Println("Creating transaction partition...")
	// Truncate to month start (simple approximation for code simplicity or use SQL)
	_, err = db.Exec(`
		DO $$
		DECLARE
			partition_date DATE := DATE_TRUNC('month', NOW());
			partition_name TEXT := 'transactions_' || TO_CHAR(partition_date, 'YYYY_MM');
			start_date DATE := partition_date;
			end_date DATE := partition_date + INTERVAL '1 month';
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = partition_name) THEN
				EXECUTE format(
					'CREATE TABLE IF NOT EXISTS %I PARTITION OF transactions FOR VALUES FROM (%L) TO (%L)',
					partition_name, start_date, end_date
				);
			END IF;
		END $$;
	`)
	if err != nil {
		log.Printf("Error creating partition: %v", err)
	}

	for i := 0; i < 10; i++ {
		hash := fmt.Sprintf("0x%d", time.Now().UnixNano()+int64(i))
		_, err := db.Exec(`
			INSERT INTO transactions (hash, from_address, to_address, amount, status, submitted_at, benchmark_id, latency_ms)
			VALUES ($1, '0xSender', '0xReceiver', 100, 'confirmed', NOW(), $2, 15.5)
			ON CONFLICT DO NOTHING`, hash, benchID)
		if err != nil {
			log.Printf("Error seeding transaction: %v", err)
		}
	}

	// 4. Metrics
	metricNames := []string{"cpu_usage", "memory_usage", "consensus_tps"}
	for _, m := range metricNames {
		_, err := db.Exec(`
			INSERT INTO metrics (timestamp, node_id, metric_name, metric_value, benchmark_id)
			VALUES (NOW(), 'node-001', $1, 45.5, $2)
			ON CONFLICT DO NOTHING`, m, benchID)
		if err != nil {
			log.Printf("Error seeding metric %s: %v", m, err)
		}
	}

	// 5. Anomalies
	_, err = db.Exec(`
		INSERT INTO anomalies (anomaly_type, severity, confidence_score, node_id, description, detected_at)
		VALUES ('sybil', 'high', 0.95, 'node-003', 'Suspected sybil behavior detected from node-003', NOW())
		ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Error seeding anomaly: %v", err)
	}
}
