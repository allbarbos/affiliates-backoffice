package middleware

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cors(t *testing.T) {
	t.Run("Should return fiber handler", func(t *testing.T) {
		got := Cors()
		assert.NotNil(t, got)
		assert.Equal(t, "func(*fiber.Ctx) error", reflect.TypeOf(got).String())
	})
}
