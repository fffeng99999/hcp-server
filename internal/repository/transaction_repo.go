package repository

import (
	"context"
	"errors"

	"github.com/fffeng99999/hcp-server/internal/models"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) GetByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	var tx models.Transaction
	if err := r.db.WithContext(ctx).Where("hash = ?", hash).First(&tx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) List(ctx context.Context, filter TransactionFilter, page, pageSize int) ([]models.Transaction, int64, error) {
	var txs []models.Transaction
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Transaction{})

	if filter.BenchmarkID != "" {
		query = query.Where("benchmark_id = ?", filter.BenchmarkID)
	}
	if filter.FromAddress != "" {
		query = query.Where("from_address = ?", filter.FromAddress)
	}
	if filter.ToAddress != "" {
		query = query.Where("to_address = ?", filter.ToAddress)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("submitted_at DESC").Offset(offset).Limit(pageSize).Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

func (r *transactionRepository) GetStats(ctx context.Context, benchmarkID string) (*TransactionStats, error) {
	var stats TransactionStats
	var result struct {
		Total        int64
		Pending      int64
		Confirmed    int64
		Failed       int64
		AvgLatencyMs float64
	}

	// This is a simplified aggregation. In production, this might need optimization or raw SQL.
	// Assuming latency_ms is populated only for confirmed transactions or valid ones.
	
	// Count totals by status
	// We can do this in one query using FILTER or CASE WHEN (Postgres)
	// Or multiple queries. For simplicity and GORM compatibility:
	
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'confirmed') as confirmed,
			COUNT(*) FILTER (WHERE status = 'failed') as failed,
			COALESCE(AVG(latency_ms), 0) as avg_latency_ms
		FROM transactions
		WHERE benchmark_id = ?
	`, benchmarkID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	stats.TotalTransactions = result.Total
	stats.PendingCount = result.Pending
	stats.ConfirmedCount = result.Confirmed
	stats.FailedCount = result.Failed
	stats.AvgLatencyMs = result.AvgLatencyMs

	return &stats, nil
}
