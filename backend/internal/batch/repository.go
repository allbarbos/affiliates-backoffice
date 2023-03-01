package batch

import (
	"affiliates-backoffice-backend/pkg/log"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	log log.LoggerI
	db  *gorm.DB
}

type RepositoryI interface {
	BeginTran() *gorm.DB
	SaveFile(ctx context.Context, model *Model) (int64, error)
	SaveErrors(ctx context.Context, id uuid.UUID, errs []FileErrorModel) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetFileToProcess(ctx context.Context) (*Model, error)
	GetFiles(ctx context.Context, affiliateID string) (*[]Model, error)
}

func NewRepository(logger log.LoggerI, db *gorm.DB) RepositoryI {
	return &repository{
		log: logger,
		db:  db,
	}
}

func (r *repository) BeginTran() *gorm.DB {
	tx := r.db.Begin()
	return tx
}

func (r *repository) SaveFile(ctx context.Context, model *Model) (int64, error) {
	tx := r.db.Create(model)
	if tx.Error != nil {
		r.log.With(ctx).Error("error saving to database: ", tx.Error)
		return 0, errors.New("error saving to database")
	}
	return tx.RowsAffected, nil
}

func (r *repository) GetFileToProcess(ctx context.Context) (*Model, error) {
	var batch Model
	tx := r.db.Where("status = 'CREATED'").First(&batch)
	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return nil, nil
		}
		r.log.With(ctx).Error(tx.Error)
		return nil, errors.New("error fetching batch from database")
	}

	tx = r.db.Model(&batch).Update("status", "PROCESSING")
	if tx.Error != nil {
		r.log.With(ctx).Error(tx.Error)
		return nil, errors.New("error updating batch status")
	}

	return &batch, nil
}

func (r *repository) GetFiles(ctx context.Context, affiliateID string) (*[]Model, error) {
	var files []Model
	tx := r.db.Where("affiliate_id = ?", affiliateID).Find(&files)
	if tx.Error != nil {
		r.log.With(ctx).Error(tx.Error)
		return nil, errors.New("error getting batch from database")
	}
	return &files, nil
}

func (r *repository) SaveErrors(ctx context.Context, id uuid.UUID, errs []FileErrorModel) error {
	sql := `UPDATE public.batches SET status='ERROR', errors=?::jsonb WHERE id=?;`
	bytes, _ := json.Marshal(errs)
	tx := r.db.Exec(sql, string(bytes), id)
	if tx.Error != nil {
		r.log.With(ctx).Error(tx.Error)
		return errors.New("error to save errors in database")
	}
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	sql := `UPDATE public.batches SET status=? WHERE id=?;`
	tx := r.db.Exec(sql, status, id)
	if tx.Error != nil {
		r.log.With(ctx).Error(tx.Error)
		return errors.New("error to save errors in database")
	}
	return nil
}
