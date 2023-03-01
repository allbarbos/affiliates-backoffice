package batch

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID          uuid.UUID  `gorm:"column:id"`
	AffiliateID uuid.UUID  `gorm:"column:affiliate_id"`
	BatchRaw    string     `gorm:"column:batch_raw"`
	Status      string     `gorm:"column:status"`
	Errors      *[]byte    `gorm:"column:errors;default:'[]'"`
	CreatedAt   *time.Time `gorm:"column:created_at;default:now"`
}

func (Model) TableName() string {
	return "batches"
}

type FileErrorModel struct {
	Row    int      `json:"row"`
	Errors []string `json:"errors,omitempty"`
}
