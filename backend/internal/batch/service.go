package batch

import (
	"affiliates-backoffice-backend/pkg/log"
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

type ServiceI interface {
	Save(ctx context.Context, affiliateID string, attachment *multipart.FileHeader) (*uuid.UUID, error)
	GetFiles(ctx context.Context, affiliateID string) (*[]Model, error)
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

func (s *service) Save(ctx context.Context, affiliateID string, attachment *multipart.FileHeader) (*uuid.UUID, error) {
	file, err := attachment.Open()
	if err != nil {
		s.log.With(ctx).Error(err)
		return nil, errors.New("error opening attachment")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.log.With(ctx).Error(err)
		return nil, errors.New("error reading attachment")
	}

	m := &Model{
		ID:          uuid.New(),
		AffiliateID: uuid.MustParse(affiliateID),
		BatchRaw:    string(fileBytes),
		Status:      "CREATED",
	}
	_, err = s.repo.SaveFile(ctx, m)
	if err != nil {
		s.log.With(ctx).Error(err)
		return nil, err
	}

	return &m.ID, nil
}

func (s *service) GetFiles(ctx context.Context, affiliateID string) (*[]Model, error) {
	return s.repo.GetFiles(ctx, affiliateID)
}
