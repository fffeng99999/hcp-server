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
	// UpdateMetrics(ctx context.Context, id string, metrics ...) error // 预留接口：后续可扩展节点指标更新
}

type nodeService struct {
	repo repository.NodeRepository
}

func NewNodeService(repo repository.NodeRepository) NodeService {
	return &nodeService{repo: repo}
}

func (s *nodeService) Register(ctx context.Context, node *models.Node) (*models.Node, error) {
	// 若节点已存在则更新信息，否则创建新节点
	existing, err := s.repo.GetByID(ctx, node.ID)
	if err == nil && existing != nil {
		// 更新已有节点记录
		// 如需更精细字段控制，可在此处补充业务逻辑
		node.RegisteredAt = existing.RegisteredAt // 保留原始注册时间
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
