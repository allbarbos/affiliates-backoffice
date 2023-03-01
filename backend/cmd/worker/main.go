package main

import (
	"affiliates-backoffice-backend/internal/batch"
	"affiliates-backoffice-backend/internal/database"
	"affiliates-backoffice-backend/internal/transaction"
	"affiliates-backoffice-backend/internal/worker"
	"affiliates-backoffice-backend/pkg/config"
	"affiliates-backoffice-backend/pkg/log"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	config.Init()
	logger := log.New()
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("logger sync() error: ", err)
		}
	}()

	database.Connect(logger)

	tranRepo := transaction.NewRepository(logger, database.DB)
	batchRepo := batch.NewRepository(logger, database.DB)
	srv := worker.NewService(logger, batchRepo, tranRepo)

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every(3).Seconds().Do(srv.ProcessBatch)
	s.StartBlocking()
}
