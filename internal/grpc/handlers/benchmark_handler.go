package handlers

import (
	"context"

	pb "github.com/fffeng99999/hcp-server/api/generated/benchmark"
	common "github.com/fffeng99999/hcp-server/api/generated/common"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/service"
)

type BenchmarkHandler struct {
	pb.UnimplementedBenchmarkServiceServer
	svc service.BenchmarkService
}

func NewBenchmarkHandler(svc service.BenchmarkService) *BenchmarkHandler {
	return &BenchmarkHandler{svc: svc}
}

func (h *BenchmarkHandler) CreateBenchmark(ctx context.Context, req *pb.CreateBenchmarkRequest) (*pb.CreateBenchmarkResponse, error) {
	benchmark := &models.Benchmark{
		Name:        req.Name,
		Description: req.Description,
		Algorithm:   req.Algorithm,
		NodeCount:   int(req.NodeCount),
		Duration:    int(req.Duration),
		TargetTPS:   int(req.TargetTps),
	}

	created, err := h.svc.Create(ctx, benchmark)
	if err != nil {
		return nil, err
	}

	return &pb.CreateBenchmarkResponse{
		Benchmark: mapModelToProto(created),
	}, nil
}

func (h *BenchmarkHandler) GetBenchmark(ctx context.Context, req *pb.GetBenchmarkRequest) (*pb.GetBenchmarkResponse, error) {
	benchmark, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetBenchmarkResponse{
		Benchmark: mapModelToProto(benchmark),
	}, nil
}

func (h *BenchmarkHandler) ListBenchmarks(ctx context.Context, req *pb.ListBenchmarksRequest) (*pb.ListBenchmarksResponse, error) {
	page := 1
	pageSize := 10
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.PageSize > 0 {
			pageSize = int(req.Pagination.PageSize)
		}
	}

	benchmarks, total, err := h.svc.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbBenchmarks []*pb.Benchmark
	for _, b := range benchmarks {
		pbBenchmarks = append(pbBenchmarks, mapModelToProto(&b))
	}

	return &pb.ListBenchmarksResponse{
		Benchmarks: pbBenchmarks,
		Pagination: &common.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  int32((total + int64(pageSize) - 1) / int64(pageSize)),
			CurrentPage: int32(page),
		},
	}, nil
}

func (h *BenchmarkHandler) UpdateBenchmark(ctx context.Context, req *pb.UpdateBenchmarkRequest) (*pb.UpdateBenchmarkResponse, error) {
	// Not implemented for skeleton
	return &pb.UpdateBenchmarkResponse{}, nil
}

func (h *BenchmarkHandler) DeleteBenchmark(ctx context.Context, req *pb.DeleteBenchmarkRequest) (*common.StatusResponse, error) {
	// Not implemented for skeleton
	return &common.StatusResponse{Success: true}, nil
}


// Helper
func mapModelToProto(m *models.Benchmark) *pb.Benchmark {
	return &pb.Benchmark{
		Id:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		Algorithm:   m.Algorithm,
		NodeCount:   int32(m.NodeCount),
		Duration:    int32(m.Duration),
		TargetTps:   int32(m.TargetTPS),
		Status:      m.Status,
		ActualTps:   m.ActualTPS,
		LatencyAvg:  m.LatencyAvg,
		CreatedAt:   m.CreatedAt.String(),
	}
}
