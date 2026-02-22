package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
	common "github.com/fffeng99999/hcp-server/api/generated/common"
	pb "github.com/fffeng99999/hcp-server/api/generated/transaction"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
	"github.com/fffeng99999/hcp-server/internal/service"
)

type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	// 校验基准测试 ID 的合法性
	benchmarkID, err := uuid.Parse(req.BenchmarkId)
	if err != nil {
		// 当前采用严格模式：解析失败直接返回错误
		// 真实场景下建议返回 InvalidArgument 类型的 gRPC 错误
		return nil, err
	}

	tx := &models.Transaction{
		Hash:        uuid.New().String(), // 目前使用随机 UUID 作为哈希；生产环境应替换为真实哈希
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Amount:      req.Amount,
		BenchmarkID: benchmarkID,
		Status:      "pending",
		SubmittedAt: time.Now(),
	}

	created, err := h.svc.Create(ctx, tx)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTransactionResponse{
		Transaction: mapTransactionToProto(created),
	}, nil
}

func (h *TransactionHandler) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	tx, err := h.svc.Get(ctx, req.Hash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil // 当前返回空结果；后续可改为返回 NotFound 错误
	}
	return &pb.GetTransactionResponse{
		Transaction: mapTransactionToProto(tx),
	}, nil
}

func (h *TransactionHandler) ListTransactions(ctx context.Context, req *pb.ListTransactionsRequest) (*pb.ListTransactionsResponse, error) {
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

	filter := repository.TransactionFilter{
		BenchmarkID: req.BenchmarkId,
		FromAddress: req.FromAddress,
		ToAddress:   req.ToAddress,
		Status:      req.Status,
	}

	txs, total, err := h.svc.List(ctx, filter, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbTxs []*pb.Transaction
	for _, tx := range txs {
		pbTxs = append(pbTxs, mapTransactionToProto(&tx))
	}

	return &pb.ListTransactionsResponse{
		Transactions: pbTxs,
		Pagination: &common.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  int32((total + int64(pageSize) - 1) / int64(pageSize)),
			CurrentPage: int32(page),
		},
	}, nil
}

func (h *TransactionHandler) GetTransactionStats(ctx context.Context, req *pb.GetTransactionStatsRequest) (*pb.GetTransactionStatsResponse, error) {
	stats, err := h.svc.GetStats(ctx, req.BenchmarkId)
	if err != nil {
		return nil, err
	}

	// 计算 TPS（当前为简化占位实现）
	tps := 0.0
	// 理想情况下应基于时间窗口统计；目前仅返回占位值

	return &pb.GetTransactionStatsResponse{
		TotalTransactions: stats.TotalTransactions,
		PendingCount:      stats.PendingCount,
		ConfirmedCount:    stats.ConfirmedCount,
		FailedCount:       stats.FailedCount,
		AvgLatencyMs:      stats.AvgLatencyMs,
		Tps:               tps,
	}, nil
}

func mapTransactionToProto(t *models.Transaction) *pb.Transaction {
	pbTx := &pb.Transaction{
		Hash:             t.Hash,
		FromAddress:      t.FromAddress,
		ToAddress:        t.ToAddress,
		Amount:           t.Amount,
		GasPrice:         t.GasPrice,
		GasLimit:         t.GasLimit,
		GasUsed:          t.GasUsed,
		Nonce:            t.Nonce,
		BlockNumber:      t.BlockNumber,
		BlockHash:        t.BlockHash,
		TransactionIndex: int32(t.TransactionIndex),
		Status:           t.Status,
		ErrorMessage:     t.ErrorMessage,
		SubmittedAt:      t.SubmittedAt.Format(time.RFC3339),
		LatencyMs:        t.LatencyMs,
		BenchmarkId:      t.BenchmarkID.String(),
	}
	if t.ConfirmedAt != nil {
		pbTx.ConfirmedAt = t.ConfirmedAt.Format(time.RFC3339)
	}
	return pbTx
}
