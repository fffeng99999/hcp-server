package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	common "github.com/fffeng99999/hcp-server/api/generated/common"
	pb "github.com/fffeng99999/hcp-server/api/generated/metric"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/service"
)

type MetricHandler struct {
	pb.UnimplementedMetricServiceServer
	svc service.MetricService
}

func NewMetricHandler(svc service.MetricService) *MetricHandler {
	return &MetricHandler{svc: svc}
}

func (h *MetricHandler) ReportMetric(ctx context.Context, req *pb.ReportMetricRequest) (*pb.ReportMetricResponse, error) {
	benchmarkID, _ := uuid.Parse(req.BenchmarkId)

	var labels map[string]interface{}
	if req.LabelsJson != "" {
		_ = json.Unmarshal([]byte(req.LabelsJson), &labels)
	}

	metric := &models.Metric{
		Timestamp:   time.Now(),
		NodeID:      req.NodeId,
		MetricName:  req.MetricName,
		MetricValue: req.MetricValue,
		MetricUnit:  req.MetricUnit,
		Labels:      labels,
		BenchmarkID: benchmarkID,
	}

	err := h.svc.Report(ctx, metric)
	if err != nil {
		return &pb.ReportMetricResponse{Success: false}, err
	}

	return &pb.ReportMetricResponse{Success: true}, nil
}

func (h *MetricHandler) GetNodeMetrics(ctx context.Context, req *pb.GetNodeMetricsRequest) (*pb.GetNodeMetricsResponse, error) {
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

	var startTime, endTime time.Time
	if req.StartTime != "" {
		startTime, _ = time.Parse(time.RFC3339, req.StartTime)
	}
	if req.EndTime != "" {
		endTime, _ = time.Parse(time.RFC3339, req.EndTime)
	}

	metrics, total, err := h.svc.GetNodeMetrics(ctx, req.NodeId, req.MetricName, startTime, endTime, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbMetrics []*pb.Metric
	for _, m := range metrics {
		pbMetrics = append(pbMetrics, mapMetricToProto(&m))
	}

	return &pb.GetNodeMetricsResponse{
		Metrics: pbMetrics,
		Pagination: &common.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  int32((total + int64(pageSize) - 1) / int64(pageSize)),
			CurrentPage: int32(page),
		},
	}, nil
}

func (h *MetricHandler) GetBenchmarkMetrics(ctx context.Context, req *pb.GetBenchmarkMetricsRequest) (*pb.GetBenchmarkMetricsResponse, error) {
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

	metrics, total, err := h.svc.GetBenchmarkMetrics(ctx, req.BenchmarkId, req.MetricName, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbMetrics []*pb.Metric
	for _, m := range metrics {
		pbMetrics = append(pbMetrics, mapMetricToProto(&m))
	}

	return &pb.GetBenchmarkMetricsResponse{
		Metrics: pbMetrics,
		Pagination: &common.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  int32((total + int64(pageSize) - 1) / int64(pageSize)),
			CurrentPage: int32(page),
		},
	}, nil
}

func mapMetricToProto(m *models.Metric) *pb.Metric {
	labelsJson, _ := json.Marshal(m.Labels)
	return &pb.Metric{
		Timestamp:    m.Timestamp.Format(time.RFC3339),
		NodeId:       m.NodeID,
		MetricName:   m.MetricName,
		MetricValue:  m.MetricValue,
		MetricUnit:   m.MetricUnit,
		LabelsJson:   string(labelsJson),
		BenchmarkId:  m.BenchmarkID.String(),
	}
}
