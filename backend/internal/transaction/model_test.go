package transaction

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ParseRowToModel(t *testing.T) {
	t.Run("Should return table name", func(t *testing.T) {
		m := Model{}
		if got := m.TableName(); got != "transactions" {
			t.Errorf("model.TableName() = %v, want %v", got, "transactions")
		}
	})

	t.Run("Should parse row to model with positive balance", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS         "

		model, errs := ParseRowToModel(batchID, affiliateID, row)

		assert.Len(t, *errs, 0)
		assert.Equal(t, batchID, model.BatchID)
		assert.Equal(t, affiliateID, model.AffiliateID)
		assert.Equal(t, 1, model.Type)
		assert.Equal(t, "Curso De Bem-Estar", model.Product)
		assert.Equal(t, 127.5, model.Value)
		assert.Equal(t, "Jose Carlos", model.Seller)
	})

	t.Run("Should parse row to model with negative balance", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "32022-02-04T07:42:12-03:00DESENVOLVEDOR FULL STACK      0000050000CELSO DE MELO       "

		model, errs := ParseRowToModel(batchID, affiliateID, row)

		assert.Len(t, *errs, 0)
		assert.Equal(t, batchID, model.BatchID)
		assert.Equal(t, affiliateID, model.AffiliateID)
		assert.Equal(t, 3, model.Type)
		assert.Equal(t, "Desenvolvedor Full Stack", model.Product)
		assert.Equal(t, -500.0, model.Value)
		assert.Equal(t, "Celso De Melo", model.Seller)
	})

	t.Run("Should return error if line is less than 86 characters", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS"

		model, errss := ParseRowToModel(batchID, affiliateID, row)
		errs := *errss

		assert.Nil(t, model)
		assert.EqualError(t, errs[0], "number of columns is different from 86")
	})

	t.Run("Should error if value is invalid", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            000001275JJOSE CARLOS         "

		_, errss := ParseRowToModel(batchID, affiliateID, row)
		errs := *errss

		assert.EqualError(t, errs[0], "the value of sale is not parsed")
	})

	t.Run("Should error if date is invalid", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "12022-01-15T19:20:30-03:CCCURSO DE BEM-ESTAR            0000012750JOSE CARLOS         "

		_, errss := ParseRowToModel(batchID, affiliateID, row)
		errs := *errss

		assert.EqualError(t, errs[0], "the date of sale is not parsed")
	})

	t.Run("Should error if type is invalid", func(t *testing.T) {
		batchID := uuid.New()
		affiliateID := uuid.New()
		row := "A2022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS         "

		_, errss := ParseRowToModel(batchID, affiliateID, row)
		errs := *errss

		assert.EqualError(t, errs[0], "the type of sale is not parsed")
		assert.EqualError(t, errs[1], "unsupported type")
	})
}
