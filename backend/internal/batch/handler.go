package batch

import (
	"affiliates-backoffice-backend/internal/errors"
	ctxpkg "affiliates-backoffice-backend/pkg/context"
	"affiliates-backoffice-backend/pkg/log"
	"context"
	"fmt"

	"net/http"

	"github.com/gofiber/fiber/v2"
)

type resource struct {
	log log.LoggerI
	srv ServiceI
}

// New instantiate a new batch handler
func NewHandler(logger log.LoggerI, service ServiceI) *resource {
	return &resource{
		log: logger,
		srv: service,
	}
}

// Get handler for route /affiliates/:affiliateID/batches
func (r *resource) Get(c *fiber.Ctx) error {
	affiliateID := c.Params("affiliateID")
	ctx := context.Background()
	ctxpkg.SetValue(&ctx, ctxpkg.KeyCorrelationID, affiliateID)
	r.log.With(ctx).Info("get batches")

	files, err := r.srv.GetFiles(ctx, affiliateID)
	if err != nil {
		r.log.Error(err)
		return errors.RenderError(c, http.StatusInternalServerError, "Error from get files")
	}

	return renderBatches(c, files)
}

// Post handler for route /affiliates/:affiliateID/batches
func (r *resource) Post(c *fiber.Ctx) error {
	affiliateID := c.Params("affiliateID")
	ctx := context.Background()
	ctxpkg.SetValue(&ctx, ctxpkg.KeyCorrelationID, affiliateID)
	r.log.With(ctx).Info("save a new batch")

	attachment, err := c.FormFile("attachment")
	if err != nil {
		r.log.With(ctx).Error(err)
		return errors.RenderError(c, http.StatusBadRequest, "Unable to receive attachment")
	}

	fileType := attachment.Header.Get("Content-Type")
	if fileType != "text/plain" {
		r.log.With(ctx).Errorf("Attachment content-type %v unsupported", fileType)
		return errors.RenderError(c, http.StatusBadRequest, fmt.Sprintf("Attachment content-type %v unsupported", fileType))
	}

	if attachment.Size == 0 {
		r.log.With(ctx).Error("Attachment cannot be empty")
		return errors.RenderError(c, http.StatusBadRequest, "Attachment cannot be empty")
	}

	batchID, err := r.srv.Save(ctx, affiliateID, attachment)
	if err != nil {
		r.log.With(ctx).Error(err)
		return errors.RenderError(c, http.StatusUnprocessableEntity, err.Error())
	}

	return renderAccepted(c, *batchID)
}
