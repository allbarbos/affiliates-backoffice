package errors

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_RenderError(t *testing.T) {
	t.Run("Should render error", func(t *testing.T) {
		app := fiber.New()
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		err := RenderError(ctx, 500, "test error")
		assert.Nil(t, err)
	})
}
