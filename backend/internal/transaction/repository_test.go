package transaction

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
	rs.model = Model{
		BatchID:     uuid.New(),
		AffiliateID: uuid.MustParse("eb2bb565-a4f1-48f1-be22-691d6ef64175"),
		Type:        1,
		Date:        time.Now(),
		Product:     "Product",
		Value:       1.40,
		Seller:      "Seller",
	}
}

func (rs *repositorySuite) AfterTest(_, _ string) {
	assert.NoError(rs.T(), rs.mockSQL.ExpectationsWereMet())
}

func (rs *repositorySuite) Test_GetTransactions() {
	rs.T().Run("Should get transactions by affiliate", func(t *testing.T) {
		const sql = `SELECT * FROM "transactions" WHERE affiliate_id = $1`
		rows := sqlmock.NewRows(
			[]string{"batch_id", "affiliate_id", "type", "date", "product", "value", "seller"},
		).
			AddRow(
				rs.model.BatchID,
				rs.model.AffiliateID,
				rs.model.Type,
				rs.model.Date,
				rs.model.Product,
				rs.model.Value,
				rs.model.Seller,
			)
		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql)).
			WithArgs(rs.model.AffiliateID).
			WillReturnRows(rows)

		items, err := rs.repo.GetTransactions(rs.model.AffiliateID.String())
		i := *items
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(rs.model, i[0]))

	})

	rs.T().Run("Should return error if database fails", func(t *testing.T) {
		const sql = `SELECT * FROM "transactions" WHERE affiliate_id = $1`
		rs.mockSQL.ExpectQuery(
			regexp.QuoteMeta(sql)).
			WithArgs(rs.model.AffiliateID).
			WillReturnError(errors.New("mock error"))

		items, err := rs.repo.GetTransactions(rs.model.AffiliateID.String())
		assert.Nil(t, items)
		assert.EqualError(t, err, "error getting transactions in database")

	})

}

// func (rs *repositorySuite) Test_Save() {
// 	rs.T().Run("Should salve transactions", func(t *testing.T) {
// 		const sql = `INSERT INTO "transactions" ("batch_id","affiliate_id","type","date","product","value","seller") VALUES ($1,$2,$3,$4,$5,$6,$7)`
// 		rs.mockSQL.ExpectBegin()
// 		rs.mockSQL.ExpectQuery(
// 			regexp.QuoteMeta(sql),
// 		).
// 			WithArgs(
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 				sqlmock.AnyArg(),
// 			)
// 		rs.mockSQL.ExpectCommit()

// 		dbTran := rs.mockDB.Begin()
// 		items := []*Model{&rs.model}

// 		err := rs.repo.Save(dbTran, items)
// 		assert.NoError(rs.T(), err)
// 	})
// }
