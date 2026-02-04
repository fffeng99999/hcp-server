package models

import (
	"time"
)

type Node struct {
	ID string `gorm:"type:varchar(50);primary_key" json:"id"`

	// Basic Info
	Name      string `gorm:"type:varchar(255)" json:"name"`
	Address   string `gorm:"type:varchar(255);not null" json:"address"`
	PublicKey string `gorm:"type:text" json:"public_key"`
	Region    string `gorm:"type:varchar(50)" json:"region"`

	// Role & Status
	Role   string `gorm:"type:varchar(20);default:'validator';index" json:"role"` // leader/validator/follower
	Status string `gorm:"type:varchar(20);default:'offline';index" json:"status"` // online/offline/syncing/failed

	// Metrics
	TrustScore           float64 `gorm:"type:decimal(5,2);default:100.00;index" json:"trust_score"`
	UptimePercentage     float64 `gorm:"type:decimal(5,2)" json:"uptime_percentage"`
	TotalBlocksProposed  int     `gorm:"default:0" json:"total_blocks_proposed"`
	TotalBlocksValidated int     `gorm:"default:0" json:"total_blocks_validated"`

	// Resource
	CPUUsage    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage"`
	MemoryUsage float64 `gorm:"type:decimal(10,2)" json:"memory_usage"`
	DiskUsage   float64 `gorm:"type:decimal(10,2)" json:"disk_usage"`

	// Network
	PeersCount        int     `gorm:"default:0" json:"peers_count"`
	NetworkLatencyAvg float64 `gorm:"type:decimal(10,4)" json:"network_latency_avg"`

	// Time
	LastHeartbeat *time.Time `gorm:"index" json:"last_heartbeat"`
	RegisteredAt  time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"registered_at"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
