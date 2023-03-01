package batch_test

import (
	"affiliates-backoffice-backend/internal/batch"
	"testing"
)

func Test_model_TableName(t *testing.T) {
	t.Run("Should return table name", func(t *testing.T) {
		m := batch.Model{}
		if got := m.TableName(); got != "batches" {
			t.Errorf("model.TableName() = %v, want %v", got, "batches")
		}
	})
}
