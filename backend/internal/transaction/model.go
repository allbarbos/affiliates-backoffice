package transaction

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Model struct {
	BatchID     uuid.UUID `gorm:"column:batch_id"`
	AffiliateID uuid.UUID `gorm:"column:affiliate_id"`
	Type        int       `gorm:"column:type"`
	Date        time.Time `gorm:"column:date"`
	Product     string    `gorm:"column:product"`
	Value       float64   `gorm:"column:value"`
	Seller      string    `gorm:"column:seller"`
}

func (Model) TableName() string {
	return "transactions"
}

func ParseRowToModel(batchID uuid.UUID, affiliateID uuid.UUID, row string) (*Model, *[]error) {
	errs := &[]error{}
	if len(row) != 86 {
		*errs = append(*errs, errors.New("number of columns is different from 86"))
		return nil, errs
	}

	tran := &Model{
		BatchID:     batchID,
		AffiliateID: affiliateID,
	}

	parseType(row, tran, errs)
	parseProduct(row, tran, errs)
	parseSeller(row, tran, errs)
	parseDate(row, tran, errs)
	parseValue(row, tran, errs)

	return tran, errs
}
