package middleware

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateHeaders(t *testing.T) {
	t.Run("Should return fiber handler", func(t *testing.T) {
		middleware := ValidateHeaders()
		assert.NotNil(t, middleware)
		assert.Equal(t, "func(*fiber.Ctx) error", reflect.TypeOf(middleware).String())

		app := fiber.New()
		app.Use(middleware)
		app.Use("/", func(c *fiber.Ctx) error {
			return c.Next()
		})
		app.Get("/test", func(c *fiber.Ctx) error {
			return nil
		})

		req := httptest.NewRequest(fiber.MethodGet, "http://without-header.com/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, 400, resp.StatusCode)

		req = httptest.NewRequest(fiber.MethodGet, "http://with-header.com/test", nil)
		req.Header.Add("X-Request-Id", "test")
		resp, _ = app.Test(req)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
