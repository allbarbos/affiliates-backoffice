package worker

import (
	"affiliates-backoffice-backend/internal/batch"
	"affiliates-backoffice-backend/internal/transaction"
	ctxpkg "affiliates-backoffice-backend/pkg/context"
	"affiliates-backoffice-backend/pkg/log"
	"context"
	"regexp"
	"sync"

	"github.com/google/uuid"
)

type ServiceI interface {
	ProcessBatch()
}

type service struct {
	log       log.LoggerI
	batchRepo batch.RepositoryI
	tranRepo  transaction.RepositoryI
}

func NewService(logger log.LoggerI, batchRepo batch.RepositoryI, tranRepo transaction.RepositoryI) ServiceI {
	return &service{
		log:       logger,
		batchRepo: batchRepo,
		tranRepo:  tranRepo,
	}
}

func (s *service) ProcessBatch() {
	ctx := context.Background()
	batchFile, err := s.batchRepo.GetFileToProcess(ctx)
	if err != nil {
		s.log.Fatal(err)
	}

	if batchFile == nil || batchFile.BatchRaw == "" {
		s.log.Info("No files to process")
		return
	}

	ctxpkg.SetValue(&ctx, ctxpkg.KeyCorrelationID, batchFile.AffiliateID.String())
	s.log.With(ctx).Info("No files to process")
	r := regexp.MustCompile(`[^\n]+`)
	rows := r.FindAllString(batchFile.BatchRaw, -1)

	transactions := make([]*transaction.Model, len(rows))
	var wg sync.WaitGroup
	errsChan := make(chan batch.FileErrorModel)
	defer close(errsChan)

	for i, row := range rows {
		wg.Add(1)
		go parallel(&wg, errsChan, transactions, batchFile.ID, batchFile.AffiliateID, i, row)
	}

	var fileErrs []batch.FileErrorModel
	go handleErrors(errsChan, &fileErrs)
	wg.Wait()

	if len(fileErrs) > 0 {
		err = s.batchRepo.SaveErrors(ctx, batchFile.ID, fileErrs)
		s.log.With(ctx).Error(err)
		return
	}

	dbTran := s.batchRepo.BeginTran()
	err = s.tranRepo.Save(dbTran, transactions)
	if err != nil {
		dbTran.Rollback()
		return
	}
	err = s.batchRepo.UpdateStatus(ctx, batchFile.ID, "PROCESSED")
	if err != nil {
		dbTran.Rollback()
		return
	}
	dbTran.Commit()
}

func parallel(wg *sync.WaitGroup, errsChan chan batch.FileErrorModel, transactions []*transaction.Model, batchID, affiliateID uuid.UUID, i int, row string) {
	defer wg.Done()
	tran, errs := transaction.ParseRowToModel(batchID, affiliateID, row)
	if len(*errs) > 0 {
		fileErr := batch.FileErrorModel{
			Row: i,
		}

		for _, err := range *errs {
			fileErr.Errors = append(fileErr.Errors, err.Error())
		}

		errsChan <- fileErr
	}
	transactions[i] = tran
}

func handleErrors(errsChan chan batch.FileErrorModel, fileErrs *[]batch.FileErrorModel) {
	for e := range errsChan {
		*fileErrs = append(*fileErrs, e)
	}
}
