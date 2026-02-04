package repository

import (
	"context"
	"time"

	"github.com/fffeng99999/hcp-server/internal/models"
	"gorm.io/gorm"
)

type metricRepository struct {
	db *gorm.DB
}

func NewMetricRepository(db *gorm.DB) MetricRepository {
	return &metricRepository{db: db}
}

func (r *metricRepository) Create(ctx context.Context, metric *models.Metric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

func (r *metricRepository) CreateBatch(ctx context.Context, metrics []*models.Metric) error {
	return r.db.WithContext(ctx).Create(metrics).Error
}

func (r *metricRepository) GetNodeMetrics(ctx context.Context, nodeID, metricName string, startTime, endTime time.Time, page, pageSize int) ([]models.Metric, int64, error) {
	var metrics []models.Metric
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Metric{}).Where("node_id = ?", nodeID)

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}
	if !startTime.IsZero() {
		query = query.Where("timestamp >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("timestamp <= ?", endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&metrics).Error; err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}

func (r *metricRepository) GetBenchmarkMetrics(ctx context.Context, benchmarkID, metricName string, page, pageSize int) ([]models.Metric, int64, error) {
	var metrics []models.Metric
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Metric{}).Where("benchmark_id = ?", benchmarkID)

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("timestamp DESC").Offset(offset).Limit(pageSize).Find(&metrics).Error; err != nil {
		return nil, 0, err
	}

	return metrics, total, nil
}
