package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Anomaly struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`

	// Info
	AnomalyType     string  `gorm:"type:varchar(50);not null;index" json:"anomaly_type"`
	Severity        string  `gorm:"type:varchar(20);not null;index" json:"severity"`
	ConfidenceScore float64 `gorm:"type:decimal(5,2);not null" json:"confidence_score"`

	// Relations
	TransactionHash string      `gorm:"type:varchar(66)" json:"transaction_hash"`
	Transaction     Transaction `gorm:"foreignKey:TransactionHash" json:"transaction,omitempty"`

	NodeID string `gorm:"type:varchar(50)" json:"node_id"`
	Node   Node   `gorm:"foreignKey:NodeID" json:"node,omitempty"`

	BenchmarkID uuid.UUID `gorm:"type:uuid;index" json:"benchmark_id"`
	Benchmark   Benchmark `gorm:"foreignKey:BenchmarkID" json:"benchmark,omitempty"`

	// Evidence
	Description string                 `gorm:"type:text" json:"description"`
	Evidence    map[string]interface{} `gorm:"serializer:json;type:jsonb;index:,type:gin" json:"evidence"`

	// Status
	Status string `gorm:"type:varchar(20);default:'new';index" json:"status"`

	// Handling
	AssignedTo      string     `gorm:"type:varchar(100)" json:"assigned_to"`
	ResolvedAt      *time.Time `json:"resolved_at"`
	ResolutionNotes string     `gorm:"type:text" json:"resolution_notes"`

	// Time
	DetectedAt time.Time `gorm:"not null;index" json:"detected_at"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (a *Anomaly) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
