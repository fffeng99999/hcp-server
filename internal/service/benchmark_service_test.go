package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBenchmarkRepository 为 repository.BenchmarkRepository 的模拟实现，用于单元测试
type MockBenchmarkRepository struct {
	mock.Mock
}

func (m *MockBenchmarkRepository) Create(ctx context.Context, benchmark *models.Benchmark) error {
	args := m.Called(ctx, benchmark)
	return args.Error(0)
}

func (m *MockBenchmarkRepository) GetByID(ctx context.Context, id string) (*models.Benchmark, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Benchmark), args.Error(1)
}

func (m *MockBenchmarkRepository) List(ctx context.Context, page, pageSize int) ([]models.Benchmark, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]models.Benchmark), args.Get(1).(int64), args.Error(2)
}

func (m *MockBenchmarkRepository) Update(ctx context.Context, benchmark *models.Benchmark) error {
	args := m.Called(ctx, benchmark)
	return args.Error(0)
}

func (m *MockBenchmarkRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestBenchmarkService_Create(t *testing.T) {
	mockRepo := new(MockBenchmarkRepository)
	svc := NewBenchmarkService(mockRepo)

	ctx := context.Background()
	benchmark := &models.Benchmark{
		Name:      "Test Benchmark",
		Algorithm: "PBFT",
		NodeCount: 4,
	}

	// 期望：调用 Create 时返回 nil 错误
	mockRepo.On("Create", ctx, benchmark).Return(nil)

	// 行为：调用服务层的 Create 方法
	created, err := svc.Create(ctx, benchmark)

	// 断言：不应返回错误，且创建结果非空
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Test Benchmark", created.Name)
	mockRepo.AssertExpectations(t)
}

func TestBenchmarkService_Get(t *testing.T) {
	mockRepo := new(MockBenchmarkRepository)
	svc := NewBenchmarkService(mockRepo)

	ctx := context.Background()
	id := uuid.New().String()
	expectedBenchmark := &models.Benchmark{
		ID:   uuid.MustParse(id),
		Name: "Test Benchmark",
	}

	// 期望：GetByID 返回预期的基准测试对象
	mockRepo.On("GetByID", ctx, id).Return(expectedBenchmark, nil)

	// 行为：调用服务层的 Get 方法
	result, err := svc.Get(ctx, id)

	// 断言：返回值与预期对象相同
	assert.NoError(t, err)
	assert.Equal(t, expectedBenchmark, result)
	mockRepo.AssertExpectations(t)
}

func TestBenchmarkService_Get_NotFound(t *testing.T) {
	mockRepo := new(MockBenchmarkRepository)
	svc := NewBenchmarkService(mockRepo)

	ctx := context.Background()
	id := uuid.New().String()

	// 期望：当仓储层返回 not found 错误时，服务层也返回错误
	mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	// 行为：调用服务层的 Get 方法
	result, err := svc.Get(ctx, id)

	// 断言：应返回错误且结果为空
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
