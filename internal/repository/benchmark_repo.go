package repository

import (
	"context"

	"github.com/fffeng99999/hcp-server/internal/models"
	"gorm.io/gorm"
)

type benchmarkRepository struct {
	db *gorm.DB
}

func NewBenchmarkRepository(db *gorm.DB) BenchmarkRepository {
	return &benchmarkRepository{db: db}
}

func (r *benchmarkRepository) Create(ctx context.Context, benchmark *models.Benchmark) error {
	return r.db.WithContext(ctx).Create(benchmark).Error
}

func (r *benchmarkRepository) GetByID(ctx context.Context, id string) (*models.Benchmark, error) {
	var benchmark models.Benchmark
	if err := r.db.WithContext(ctx).First(&benchmark, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &benchmark, nil
}

func (r *benchmarkRepository) List(ctx context.Context, page, pageSize int) ([]models.Benchmark, int64, error) {
	var benchmarks []models.Benchmark
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&models.Benchmark{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&benchmarks).Error; err != nil {
		return nil, 0, err
	}

	return benchmarks, total, nil
}

func (r *benchmarkRepository) Update(ctx context.Context, benchmark *models.Benchmark) error {
	return r.db.WithContext(ctx).Save(benchmark).Error
}

func (r *benchmarkRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Benchmark{}, "id = ?", id).Error
}
