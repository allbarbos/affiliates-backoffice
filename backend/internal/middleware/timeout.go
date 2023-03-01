package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func Timeout(t time.Duration, tErrs ...error) fiber.Handler {
	h := func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
	return timeout.New(h, t, tErrs...)
}
