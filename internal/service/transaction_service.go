package service

import (
	"context"

	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
)

type TransactionService interface {
	Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error)
	Get(ctx context.Context, hash string) (*models.Transaction, error)
	List(ctx context.Context, filter repository.TransactionFilter, page, pageSize int) ([]models.Transaction, int64, error)
	GetStats(ctx context.Context, benchmarkID string) (*repository.TransactionStats, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) Create(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	// Add business logic here if needed (e.g. validation, event publishing)
	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *transactionService) Get(ctx context.Context, hash string) (*models.Transaction, error) {
	return s.repo.GetByHash(ctx, hash)
}

func (s *transactionService) List(ctx context.Context, filter repository.TransactionFilter, page, pageSize int) ([]models.Transaction, int64, error) {
	return s.repo.List(ctx, filter, page, pageSize)
}

func (s *transactionService) GetStats(ctx context.Context, benchmarkID string) (*repository.TransactionStats, error) {
	return s.repo.GetStats(ctx, benchmarkID)
}
