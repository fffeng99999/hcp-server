package repository

import (
	"context"
	"errors"

	"github.com/fffeng99999/hcp-server/internal/models"
	"gorm.io/gorm"
)

type nodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepository{db: db}
}

func (r *nodeRepository) Create(ctx context.Context, node *models.Node) error {
	return r.db.WithContext(ctx).Create(node).Error
}

func (r *nodeRepository) GetByID(ctx context.Context, id string) (*models.Node, error) {
	var node models.Node
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&node).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepository) List(ctx context.Context, filter NodeFilter, page, pageSize int) ([]models.Node, int64, error) {
	var nodes []models.Node
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Node{})

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Region != "" {
		query = query.Where("region = ?", filter.Region)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&nodes).Error; err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

func (r *nodeRepository) Update(ctx context.Context, node *models.Node) error {
	return r.db.WithContext(ctx).Save(node).Error
}

func (r *nodeRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&models.Node{}).Where("id = ?", id).Update("status", status).Error
}
