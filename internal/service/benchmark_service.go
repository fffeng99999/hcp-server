package service

import (
	"context"

	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
)

type BenchmarkService interface {
	Create(ctx context.Context, req *models.Benchmark) (*models.Benchmark, error)
	Get(ctx context.Context, id string) (*models.Benchmark, error)
	List(ctx context.Context, page, pageSize int) ([]models.Benchmark, int64, error)
}

type benchmarkService struct {
	repo repository.BenchmarkRepository
}

func NewBenchmarkService(repo repository.BenchmarkRepository) BenchmarkService {
	return &benchmarkService{repo: repo}
}

func (s *benchmarkService) Create(ctx context.Context, req *models.Benchmark) (*models.Benchmark, error) {
	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *benchmarkService) Get(ctx context.Context, id string) (*models.Benchmark, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *benchmarkService) List(ctx context.Context, page, pageSize int) ([]models.Benchmark, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}
