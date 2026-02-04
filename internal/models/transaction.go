package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Hash string `gorm:"type:varchar(66);primary_key" json:"hash"`

	// Basic Info
	FromAddress string `gorm:"type:varchar(42);not null;index" json:"from_address"`
	ToAddress   string `gorm:"type:varchar(42);not null;index" json:"to_address"`
	Amount      int64  `gorm:"not null" json:"amount"` // wei
	GasPrice    int64  `json:"gas_price"`
	GasLimit    int64  `json:"gas_limit"`
	GasUsed     int64  `json:"gas_used"`
	Nonce       int64  `json:"nonce"`

	// Block Info
	BlockNumber      int64  `gorm:"index" json:"block_number"`
	BlockHash        string `gorm:"type:varchar(66)" json:"block_hash"`
	TransactionIndex int    `json:"transaction_index"`

	// Status
	Status       string `gorm:"type:varchar(20);not null;index" json:"status"` // pending/confirmed/failed
	ErrorMessage string `gorm:"type:text" json:"error_message"`

	// Time
	SubmittedAt time.Time  `gorm:"not null;index" json:"submitted_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`

	// Metrics
	LatencyMs float64 `gorm:"type:decimal(10,4)" json:"latency_ms"`

	// Benchmark Relation
	BenchmarkID uuid.UUID `gorm:"type:uuid;index" json:"benchmark_id"`
	Benchmark   Benchmark `gorm:"foreignKey:BenchmarkID;constraint:OnDelete:CASCADE" json:"benchmark,omitempty"`

	// Metadata
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
