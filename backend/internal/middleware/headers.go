package middleware

import (
	"net/http"

	"affiliates-backoffice-backend/internal/errors"

	"github.com/gofiber/fiber/v2"
)

func ValidateHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		if headers["X-Request-Id"] == "" {
			return errors.RenderError(c, http.StatusBadRequest, "X-Request-ID header not found")
		}
		return c.Next()
	}
}
