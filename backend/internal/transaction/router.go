package transaction

import (
	"affiliates-backoffice-backend/pkg/log"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandlers(router fiber.Router, logger log.LoggerI, service ServiceI) {
	v1 := router.Group("/v1")
	resource := NewHandler(logger, service)
	v1.Get("/affiliates/:affiliateID/transactions", resource.Get)
}
