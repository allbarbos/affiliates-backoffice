package transaction

import (
	"affiliates-backoffice-backend/internal/errors"
	ctxpkg "affiliates-backoffice-backend/pkg/context"
	"affiliates-backoffice-backend/pkg/log"
	"context"

	"net/http"

	"github.com/gofiber/fiber/v2"
)

type resource struct {
	log log.LoggerI
	srv ServiceI
}

// New instantiate a new transaction handler
func NewHandler(logger log.LoggerI, service ServiceI) *resource {
	return &resource{
		log: logger,
		srv: service,
	}
}

// Get handler for route /affiliates/:affiliateID/transactions
func (r *resource) Get(c *fiber.Ctx) error {
	affiliateID := c.Params("affiliateID")
	ctx := context.Background()
	ctxpkg.SetValue(&ctx, ctxpkg.KeyCorrelationID, affiliateID)

	r.log.With(ctx).Infof("get transactions affiliate ID %v", affiliateID)

	trans, err := r.srv.Get(ctx, affiliateID)
	if err != nil {
		r.log.With(ctx).Error(err)
		return errors.RenderError(c, http.StatusInternalServerError, err.Error())
	}
	return renderTransactions(c, *trans)
}
