package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	Timestamp  time.Time `gorm:"primaryKey;not null;index" json:"timestamp"`
	NodeID     string    `gorm:"primaryKey;type:varchar(50);index" json:"node_id"`
	MetricName string    `gorm:"primaryKey;type:varchar(100);not null;index" json:"metric_name"`

	// 关联关系
	Node        Node      `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE" json:"node,omitempty"`
	BenchmarkID uuid.UUID `gorm:"type:uuid;index" json:"benchmark_id"`
	Benchmark   Benchmark `gorm:"foreignKey:BenchmarkID;constraint:OnDelete:CASCADE" json:"benchmark,omitempty"`

	// 指标数值
	MetricValue float64 `gorm:"type:decimal(20,6);not null" json:"metric_value"`
	MetricUnit  string  `gorm:"type:varchar(20)" json:"metric_unit"`

	// 维度标签信息
	Labels map[string]interface{} `gorm:"serializer:json;type:jsonb;index:,type:gin" json:"labels"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
