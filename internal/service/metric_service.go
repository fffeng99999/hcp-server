package service

import (
	"context"
	"time"

	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
)

type MetricService interface {
	Report(ctx context.Context, metric *models.Metric) error
	ReportBatch(ctx context.Context, metrics []*models.Metric) error
	GetNodeMetrics(ctx context.Context, nodeID, metricName string, startTime, endTime time.Time, page, pageSize int) ([]models.Metric, int64, error)
	GetBenchmarkMetrics(ctx context.Context, benchmarkID, metricName string, page, pageSize int) ([]models.Metric, int64, error)
}

type metricService struct {
	repo repository.MetricRepository
}

func NewMetricService(repo repository.MetricRepository) MetricService {
	return &metricService{repo: repo}
}

func (s *metricService) Report(ctx context.Context, metric *models.Metric) error {
	return s.repo.Create(ctx, metric)
}

func (s *metricService) ReportBatch(ctx context.Context, metrics []*models.Metric) error {
	return s.repo.CreateBatch(ctx, metrics)
}

func (s *metricService) GetNodeMetrics(ctx context.Context, nodeID, metricName string, startTime, endTime time.Time, page, pageSize int) ([]models.Metric, int64, error) {
	return s.repo.GetNodeMetrics(ctx, nodeID, metricName, startTime, endTime, page, pageSize)
}

func (s *metricService) GetBenchmarkMetrics(ctx context.Context, benchmarkID, metricName string, page, pageSize int) ([]models.Metric, int64, error) {
	return s.repo.GetBenchmarkMetrics(ctx, benchmarkID, metricName, page, pageSize)
}
