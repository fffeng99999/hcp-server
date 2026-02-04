package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Benchmark struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Basic Info
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Algorithm   string `gorm:"type:varchar(50);not null" json:"algorithm"` // tPBFT/Raft/HotStuff
	NodeCount   int    `gorm:"not null" json:"node_count"`
	Duration    int    `gorm:"not null" json:"duration"` // seconds
	TargetTPS   int    `json:"target_tps"`

	// Performance Metrics
	ActualTPS   float64 `gorm:"type:decimal(10,2)" json:"actual_tps"`
	LatencyP50  float64 `gorm:"type:decimal(10,4)" json:"latency_p50"`
	LatencyP90  float64 `gorm:"type:decimal(10,4)" json:"latency_p90"`
	LatencyP99  float64 `gorm:"type:decimal(10,4)" json:"latency_p99"`
	LatencyP999 float64 `gorm:"type:decimal(10,4)" json:"latency_p999"`
	LatencyAvg  float64 `gorm:"type:decimal(10,4)" json:"latency_avg"`
	LatencyMax  float64 `gorm:"type:decimal(10,4)" json:"latency_max"`
	LatencyMin  float64 `gorm:"type:decimal(10,4)" json:"latency_min"`

	// Blockchain Metrics
	BlockCount           int     `json:"block_count"`
	TransactionCount     int     `json:"transaction_count"`
	SuccessfulTx         int     `json:"successful_tx"`
	FailedTx             int     `json:"failed_tx"`
	BlockSizeAvg         float64 `gorm:"type:decimal(10,2)" json:"block_size_avg"`
	BlockPropagationTime float64 `gorm:"type:decimal(10,4)" json:"block_propagation_time"`

	// Resource Usage
	CPUUsageAvg    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage_avg"`
	CPUUsageMax    float64 `gorm:"type:decimal(5,2)" json:"cpu_usage_max"`
	MemoryUsageAvg float64 `gorm:"type:decimal(10,2)" json:"memory_usage_avg"`
	MemoryUsageMax float64 `gorm:"type:decimal(10,2)" json:"memory_usage_max"`
	NetworkInMbps  float64 `gorm:"type:decimal(10,2)" json:"network_in_mbps"`
	NetworkOutMbps float64 `gorm:"type:decimal(10,2)" json:"network_out_mbps"`
	DiskIORead     float64 `gorm:"type:decimal(10,2)" json:"disk_io_read"`
	DiskIOWrite    float64 `gorm:"type:decimal(10,2)" json:"disk_io_write"`

	// Consensus Specific
	ViewChangeCount     int     `json:"view_change_count"`
	PreparePhaseLatency float64 `gorm:"type:decimal(10,4)" json:"prepare_phase_latency"`
	CommitPhaseLatency  float64 `gorm:"type:decimal(10,4)" json:"commit_phase_latency"`

	// Status & Metadata
	Status       string     `gorm:"type:varchar(20);default:'running'" json:"status"`
	ErrorMessage string     `gorm:"type:text" json:"error_message"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (b *Benchmark) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
