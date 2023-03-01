package middleware

import (
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Timeout(t *testing.T) {
	t.Run("Should return fiber handler", func(t *testing.T) {
		middleware := Timeout(time.Millisecond * 1)
		assert.NotNil(t, middleware)
		assert.Equal(t, "func(*fiber.Ctx) error", reflect.TypeOf(middleware).String())

		app := fiber.New()
		app.Use(middleware)
		app.Use("/", func(c *fiber.Ctx) error {
			time.Sleep(time.Minute * 5)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodGet, "http://without-header.com/test", nil)
		resp, err := app.Test(req)
		assert.Nil(t, resp)
		assert.Equal(t, "test: timeout error 1000ms", err.Error())
	})
}
