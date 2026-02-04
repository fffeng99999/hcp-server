package handlers

import (
	"context"
	"time"

	common "github.com/fffeng99999/hcp-server/api/generated/common"
	pb "github.com/fffeng99999/hcp-server/api/generated/node"
	"github.com/fffeng99999/hcp-server/internal/models"
	"github.com/fffeng99999/hcp-server/internal/repository"
	"github.com/fffeng99999/hcp-server/internal/service"
)

type NodeHandler struct {
	pb.UnimplementedNodeServiceServer
	svc service.NodeService
}

func NewNodeHandler(svc service.NodeService) *NodeHandler {
	return &NodeHandler{svc: svc}
}

func (h *NodeHandler) RegisterNode(ctx context.Context, req *pb.RegisterNodeRequest) (*pb.RegisterNodeResponse, error) {
	node := &models.Node{
		ID:        req.Address, // Use address as ID or generate new UUID? Assuming address is unique ID for now.
		Name:      req.Name,
		Address:   req.Address,
		PublicKey: req.PublicKey,
		Region:    req.Region,
		Role:      req.Role,
		Status:    "online",
		RegisteredAt: time.Now(),
		UpdatedAt:    time.Now(),
	}
	// If ID is not address, generate one. But typically p2p nodes have IDs.
	if node.ID == "" {
		node.ID = req.Address // Fallback
	}

	registered, err := h.svc.Register(ctx, node)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterNodeResponse{
		Node: mapNodeToProto(registered),
	}, nil
}

func (h *NodeHandler) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	node, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	return &pb.GetNodeResponse{
		Node: mapNodeToProto(node),
	}, nil
}

func (h *NodeHandler) UpdateNodeStatus(ctx context.Context, req *pb.UpdateNodeStatusRequest) (*pb.UpdateNodeStatusResponse, error) {
	// This might need more robust update logic in service, passing struct instead of just status
	// But service interface has UpdateStatus(id, status).
	// To update metrics, we might need another method.
	
	// For now, just update status.
	err := h.svc.UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}
	
	// Fetch updated node
	node, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	return &pb.UpdateNodeStatusResponse{
		Node: mapNodeToProto(node),
	}, nil
}

func (h *NodeHandler) ListNodes(ctx context.Context, req *pb.ListNodesRequest) (*pb.ListNodesResponse, error) {
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

	filter := repository.NodeFilter{
		Role:   req.Role,
		Status: req.Status,
		Region: req.Region,
	}

	nodes, total, err := h.svc.List(ctx, filter, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbNodes []*pb.Node
	for _, n := range nodes {
		pbNodes = append(pbNodes, mapNodeToProto(&n))
	}

	return &pb.ListNodesResponse{
		Nodes: pbNodes,
		Pagination: &common.PaginationResponse{
			TotalItems:  int32(total),
			TotalPages:  int32((total + int64(pageSize) - 1) / int64(pageSize)),
			CurrentPage: int32(page),
		},
	}, nil
}

func (h *NodeHandler) GetNetworkTopology(ctx context.Context, req *pb.GetNetworkTopologyRequest) (*pb.GetNetworkTopologyResponse, error) {
	// Just return all nodes for now
	nodes, _, err := h.svc.List(ctx, repository.NodeFilter{}, 1, 1000)
	if err != nil {
		return nil, err
	}

	var pbNodes []*pb.Node
	for _, n := range nodes {
		pbNodes = append(pbNodes, mapNodeToProto(&n))
	}

	return &pb.GetNetworkTopologyResponse{
		Nodes: pbNodes,
	}, nil
}

func mapNodeToProto(n *models.Node) *pb.Node {
	pbNode := &pb.Node{
		Id:                   n.ID,
		Name:                 n.Name,
		Address:              n.Address,
		PublicKey:            n.PublicKey,
		Region:               n.Region,
		Role:                 n.Role,
		Status:               n.Status,
		TrustScore:           n.TrustScore,
		UptimePercentage:     n.UptimePercentage,
		TotalBlocksProposed:  int32(n.TotalBlocksProposed),
		TotalBlocksValidated: int32(n.TotalBlocksValidated),
		CpuUsage:             n.CPUUsage,
		MemoryUsage:          n.MemoryUsage,
		DiskUsage:            n.DiskUsage,
		PeersCount:           int32(n.PeersCount),
		NetworkLatencyAvg:    n.NetworkLatencyAvg,
		RegisteredAt:         n.RegisteredAt.Format(time.RFC3339),
		UpdatedAt:            n.UpdatedAt.Format(time.RFC3339),
	}
	if n.LastHeartbeat != nil {
		pbNode.LastHeartbeat = n.LastHeartbeat.Format(time.RFC3339)
	}
	return pbNode
}
