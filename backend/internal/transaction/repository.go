package transaction

import (
	"affiliates-backoffice-backend/pkg/log"
	"errors"

	"gorm.io/gorm"
)

type repository struct {
	log log.LoggerI
	db  *gorm.DB
}

type RepositoryI interface {
	GetTransactions(affiliateID string) (*[]Model, error)
	Save(tran *gorm.DB, transactions []*Model) error
}

func NewRepository(logger log.LoggerI, db *gorm.DB) RepositoryI {
	return &repository{
		log: logger,
		db:  db,
	}
}

func (r *repository) GetTransactions(affiliateID string) (*[]Model, error) {
	var transactions []Model
	tx := r.db.Where("affiliate_id = ?", affiliateID).Find(&transactions)
	if tx.Error != nil {
		r.log.Error(tx.Error)
		return nil, errors.New("error getting transactions in database")
	}
	return &transactions, nil
}

func (r *repository) Save(dbTran *gorm.DB, transactions []*Model) error {
	tx := dbTran.Create(transactions)
	if tx.Error != nil {
		r.log.Error(tx.Error)
		return errors.New("error saving transactions in database")
	}
	return nil
}
