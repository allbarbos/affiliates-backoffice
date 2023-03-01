package transaction

import (
	"affiliates-backoffice-backend/pkg/log"
	"context"
)

type ServiceI interface {
	Get(ctx context.Context, affiliateID string) (*[]Model, error)
}

type service struct {
	log  log.LoggerI
	repo RepositoryI
}

func NewService(logger log.LoggerI, repo RepositoryI) ServiceI {
	return &service{
		log:  logger,
		repo: repo,
	}
}

func (s *service) Get(ctx context.Context, affiliateID string) (*[]Model, error) {
	return s.repo.GetTransactions(affiliateID)
}
