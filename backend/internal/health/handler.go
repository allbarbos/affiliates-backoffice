package health

import (
	"github.com/gofiber/fiber/v2"
)

type healthResponse struct {
	Status string `json:"status"`
}

func RegisterHandlers(r fiber.Router) {
	r.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(healthResponse{Status: "UP"})
	})
}
