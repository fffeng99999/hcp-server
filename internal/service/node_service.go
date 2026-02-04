package service

import (
	"context"

	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
)

type NodeService interface {
	Register(ctx context.Context, node *models.Node) (*models.Node, error)
	Get(ctx context.Context, id string) (*models.Node, error)
	List(ctx context.Context, filter repository.NodeFilter, page, pageSize int) ([]models.Node, int64, error)
	UpdateStatus(ctx context.Context, id, status string) error
	// UpdateMetrics(ctx context.Context, id string, metrics ...) error // Can be expanded
}

type nodeService struct {
	repo repository.NodeRepository
}

func NewNodeService(repo repository.NodeRepository) NodeService {
	return &nodeService{repo: repo}
}

func (s *nodeService) Register(ctx context.Context, node *models.Node) (*models.Node, error) {
	// Check if node exists, if so update it, else create
	existing, err := s.repo.GetByID(ctx, node.ID)
	if err == nil && existing != nil {
		// Update existing
		// Logic to update fields if necessary
		// For now, let's assume registration might update info
		node.RegisteredAt = existing.RegisteredAt // Preserve creation time
		if err := s.repo.Update(ctx, node); err != nil {
			return nil, err
		}
		return node, nil
	}

	if err := s.repo.Create(ctx, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (s *nodeService) Get(ctx context.Context, id string) (*models.Node, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *nodeService) List(ctx context.Context, filter repository.NodeFilter, page, pageSize int) ([]models.Node, int64, error) {
	return s.repo.List(ctx, filter, page, pageSize)
}

func (s *nodeService) UpdateStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
