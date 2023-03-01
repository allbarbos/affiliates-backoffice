package errors

import "github.com/gofiber/fiber/v2"

type errorsResponse struct {
	Errors []string `json:"errors"`
}

func RenderError(c *fiber.Ctx, status int, errors ...string) error {
	return c.Status(status).JSON(errorsResponse{
		Errors: errors,
	})
}
