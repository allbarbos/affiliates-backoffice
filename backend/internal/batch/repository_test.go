package batch

import (
	"affiliates-backoffice-backend/pkg/log"
	"affiliates-backoffice-backend/pkg/tests"
	"context"
	"errors"
	"reflect"
	"regexp"
	"time"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type repositorySuite struct {
	suite.Suite
	mockDB  *gorm.DB
	mockSQL sqlmock.Sqlmock
	repo    RepositoryI
	ctx     context.Context
	logger  log.LoggerI
	model   Model
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

func (rs *repositorySuite) SetupSuite() {
	rs.mockSQL, rs.mockDB = tests.BuildMockDB(rs.T())
	loggerTest, _ := log.NewForTest()
	rs.logger = loggerTest
}

func (rs *repositorySuite) BeforeTest(suiteName, testName string) {
	rs.ctx = context.Background()
	rs.repo = NewRepository(rs.logger, rs.mockDB)

	affiliateID := uuid.MustParse("eb2bb565-a4f1-48f1-be22-691d6ef64175")
	raw := `12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS        \n12021-12-03T11:46:02-03:00DOMINANDO INVESTIMENTOS       0000050000MARIA CANDIDA      \n`
	errs := []byte("[]")
	now := time.Now()
	rs.model = Model{
		ID:          uuid.New(),
		AffiliateID: affiliateID,
		BatchRaw:    string(raw),
		Status:      "CREATED",
		Errors:      &errs,
		CreatedAt:   &now,
	}
}

func (rs *repositorySuite) AfterTest(_, _ string) {
	assert.NoError(rs.T(), rs.mockSQL.ExpectationsWereMet())
}

func (rs *repositorySuite) Test_BeginTran() {
	rs.T().Run("Should return a db transaction", func(t *testing.T) {
		tran := rs.repo.BeginTran()
		assert.NotNil(rs.T(), tran)
	})
}

func (rs *repositorySuite) Test_SaveFile() {
	rs.T().Run("Should save a batch file", func(t *testing.T) {
		const sql = `INSERT INTO "batches" ("id","affiliate_id","batch_raw","status") VALUES ($1,$2,$3,$4) RETURNING "errors","created_at"`
		rs.mockSQL.ExpectBegin()
		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				sqlmock.AnyArg(),
				rs.model.AffiliateID,
				rs.model.BatchRaw,
				"CREATED",
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"errors"}).AddRow("[]"),
				sqlmock.NewRows([]string{"created_at"}).AddRow("2023-02-28 13:11:58.828"),
			)
		rs.mockSQL.ExpectCommit()

		rs.model.Errors = nil
		rs.model.CreatedAt = nil
		rowsAffected, err := rs.repo.SaveFile(rs.ctx, &rs.model)

		assert.NoError(rs.T(), err)
		assert.Equal(rs.T(), int64(1), rowsAffected)
	})

	rs.T().Run("Should return error if not saving", func(t *testing.T) {
		const sql = `INSERT INTO "batches" ("id","affiliate_id","batch_raw","status") VALUES ($1,$2,$3,$4) RETURNING "errors","created_at"`
		rs.mockSQL.ExpectBegin()
		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				sqlmock.AnyArg(),
				rs.model.AffiliateID,
				rs.model.BatchRaw,
				"CREATED",
			).
			WillReturnError(errors.New("mock error"))
		rs.mockSQL.ExpectRollback()

		rs.model.Errors = nil
		rs.model.CreatedAt = nil
		rowsAffected, err := rs.repo.SaveFile(rs.ctx, &rs.model)

		assert.EqualError(rs.T(), err, "error saving to database")
		assert.Equal(rs.T(), int64(0), rowsAffected)
	})
}

func (rs *repositorySuite) Test_GetFiles() {
	rs.T().Run("Should get all batch files", func(t *testing.T) {
		sql := `SELECT * FROM "batches" WHERE affiliate_id = $1`
		rows := sqlmock.NewRows(
			[]string{"id", "affiliate_id", "batch_raw", "status", "errors", "created_at"},
		).
			AddRow(
				rs.model.ID,
				rs.model.AffiliateID,
				rs.model.BatchRaw,
				rs.model.Status,
				rs.model.Errors,
				rs.model.CreatedAt,
			)

		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql)).
			WithArgs(rs.model.AffiliateID).
			WillReturnRows(rows)

		files, err := rs.repo.GetFiles(rs.ctx, rs.model.AffiliateID.String())
		f := *files

		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(rs.model, f[0]))
	})

	rs.T().Run("Should return error", func(t *testing.T) {
		sql := `SELECT * FROM "batches" WHERE affiliate_id = $1`
		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql)).
			WithArgs(rs.model.AffiliateID).
			WillReturnError(errors.New("mock error"))

		files, err := rs.repo.GetFiles(rs.ctx, rs.model.AffiliateID.String())

		assert.Nil(t, files)
		assert.EqualError(t, err, "error getting batch from database")
	})
}

func (rs *repositorySuite) Test_SaveErrors() {
	rs.T().Run("Should save errors", func(t *testing.T) {
		id := uuid.MustParse("0cec8085-d6d0-43f6-8fcb-3548e5203ff9")
		const sql = `UPDATE public.batches SET status='ERROR', errors=$1::jsonb WHERE id=$2;`
		rs.mockSQL.ExpectExec(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				`[{"row":1,"errors":["err msg"]}]`,
				id.String(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		errs := []FileErrorModel{{Row: 1, Errors: []string{"err msg"}}}
		err := rs.repo.SaveErrors(rs.ctx, id, errs)
		assert.NoError(rs.T(), err)
	})

	rs.T().Run("Should return error if not saving errors", func(t *testing.T) {
		id := uuid.MustParse("0cec8085-d6d0-43f6-8fcb-3548e5203ff9")
		const sql = `UPDATE public.batches SET status='ERROR', errors=$1::jsonb WHERE id=$2;`
		rs.mockSQL.ExpectExec(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				`[{"row":1,"errors":["err msg"]}]`,
				id.String(),
			).
			WillReturnError(errors.New("mock error"))

		errs := []FileErrorModel{{Row: 1, Errors: []string{"err msg"}}}
		err := rs.repo.SaveErrors(rs.ctx, id, errs)
		assert.EqualError(rs.T(), err, "error to save errors in database")
	})
}

func (rs *repositorySuite) Test_UpdateStatus() {
	rs.T().Run("Should update status", func(t *testing.T) {
		id := uuid.MustParse("0cec8085-d6d0-43f6-8fcb-3548e5203ff9")
		status := "PROCESSING"
		const sql = `UPDATE public.batches SET status=$1 WHERE id=$2;`
		rs.mockSQL.ExpectExec(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				status,
				id.String(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := rs.repo.UpdateStatus(rs.ctx, id, status)
		assert.NoError(rs.T(), err)
	})

	rs.T().Run("Should return error if not updating status", func(t *testing.T) {
		id := uuid.MustParse("0cec8085-d6d0-43f6-8fcb-3548e5203ff9")
		status := "PROCESSING"
		const sql = `UPDATE public.batches SET status=$1 WHERE id=$2;`
		rs.mockSQL.ExpectExec(
			regexp.QuoteMeta(sql),
		).
			WithArgs(
				status,
				id.String(),
			).
			WillReturnError(errors.New("mock error"))

		err := rs.repo.UpdateStatus(rs.ctx, id, status)
		assert.EqualError(rs.T(), err, "error to save errors in database")
	})
}
