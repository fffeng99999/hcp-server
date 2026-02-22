package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Hash string `gorm:"type:varchar(66);primary_key" json:"hash"`

	// 基础信息
	FromAddress string `gorm:"type:varchar(42);not null;index" json:"from_address"`
	ToAddress   string `gorm:"type:varchar(42);not null;index" json:"to_address"`
	Amount      int64  `gorm:"not null" json:"amount"` // 单位为 wei
	GasPrice    int64  `json:"gas_price"`
	GasLimit    int64  `json:"gas_limit"`
	GasUsed     int64  `json:"gas_used"`
	Nonce       int64  `json:"nonce"`

	// 区块相关信息
	BlockNumber      int64  `gorm:"index" json:"block_number"`
	BlockHash        string `gorm:"type:varchar(66)" json:"block_hash"`
	TransactionIndex int    `json:"transaction_index"`

	// 交易状态
	Status       string `gorm:"type:varchar(20);not null;index" json:"status"` // pending/confirmed/failed
	ErrorMessage string `gorm:"type:text" json:"error_message"`

	// 时间信息
	SubmittedAt time.Time  `gorm:"not null;index" json:"submitted_at"`
	ConfirmedAt *time.Time `json:"confirmed_at"`

	// 性能指标
	LatencyMs float64 `gorm:"type:decimal(10,4)" json:"latency_ms"`

	// 关联的基准测试任务
	BenchmarkID uuid.UUID `gorm:"type:uuid;index" json:"benchmark_id"`
	Benchmark   Benchmark `gorm:"foreignKey:BenchmarkID;constraint:OnDelete:CASCADE" json:"benchmark,omitempty"`

	// 元信息（记录创建时间等）
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
