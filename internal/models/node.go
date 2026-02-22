package models

import (
	"time"
)

type Node struct {
	ID string `gorm:"type:varchar(50);primary_key" json:"id"`

	// 基础信息
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Address   string `gorm:"type:varchar(255);not null" json:"address"`
	PublicKey string `gorm:"type:text" json:"public_key"`
	Region    string `gorm:"type:varchar(50)" json:"region"`

	// 节点角色与运行状态
	Role   string `gorm:"type:varchar(20);default:'validator';index" json:"role"` // leader/validator/follower
	Status string `gorm:"type:varchar(20);default:'offline';index" json:"status"` // online/offline/syncing/failed

	// 运行统计与可信度指标
	TrustScore           float64 `gorm:"type:decimal(5,2);default:100.00;index" json:"trust_score"`
	UptimePercentage     float64 `gorm:"type:decimal(5,2)" json:"uptime_percentage"`
	TotalBlocksProposed  int     `gorm:"default:0" json:"total_blocks_proposed"`
	TotalBlocksValidated int     `gorm:"default:0" json:"total_blocks_validated"`

	// 资源使用情况
	CPUUsage    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage float64 `gorm:"type:decimal(10,2)" json:"memory_usage"`
	DiskUsage   float64 `gorm:"type:decimal(10,2)" json:"disk_usage"`

	// 网络相关指标
	PeersCount        int     `gorm:"default:0" json:"peers_count"`
	NetworkLatencyAvg float64 `gorm:"type:decimal(10,4)" json:"network_latency_avg"`

	// 时间相关信息
	LastHeartbeat *time.Time `gorm:"index" json:"last_heartbeat"`
	RegisteredAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"registered_at"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
