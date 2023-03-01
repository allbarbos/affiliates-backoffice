package health

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Health(t *testing.T) {
	app := fiber.New()
	RegisterHandlers(app)

	req := httptest.NewRequest("GET", "/health", nil)

	resp, _ := app.Test(req, -1)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"status":"UP"}`, string(body))
}
