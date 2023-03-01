package transaction_test

import (
	"affiliates-backoffice-backend/internal/transaction"
	mocks "affiliates-backoffice-backend/internal/transaction/mocks"
	"affiliates-backoffice-backend/pkg/log"
	"affiliates-backoffice-backend/pkg/tests"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type handlerSuite struct {
	suite.Suite
	ctx          *fiber.Ctx
	service      *mocks.MockServiceI
	logger       log.LoggerI
	affiliateID  string
	transactions []transaction.Model
}

func (suite *handlerSuite) SetupSuite() {
	loggerTest, _ := log.NewForTest()
	suite.logger = loggerTest
}

func (suite *handlerSuite) BeforeTest(suiteName, testName string) {
	app := fiber.New()
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	suite.ctx = c

	suite.affiliateID = "ca226967-6b4a-4461-bf16-ab0108f87380"
	ctrl := gomock.NewController(suite.T())
	suite.service = mocks.NewMockServiceI(ctrl)

	dateExpected := "2023-02-26T22:00:16.616853-03:00"
	layout := "2006-01-02T15:04:05-07:00"
	d, err := time.Parse(layout, dateExpected)
	assert.NoError(suite.T(), err)
	suite.transactions = []transaction.Model{
		{
			BatchID:     uuid.MustParse("4b1362ba-b443-43b4-b057-0d5431a9bbcb"),
			AffiliateID: uuid.MustParse(suite.affiliateID),
			Product:     "CURSO DE BEM-ESTAR",
			Type:        1,
			Date:        d,
			Value:       1.30,
			Seller:      "JOSE CARLOS",
		},
	}
}

func Test_HandlerSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (suite *handlerSuite) Test_Get() {
	suite.T().Run("Should ...", func(t *testing.T) {
		suite.service.EXPECT().Get(gomock.Any(), suite.affiliateID).Return(&suite.transactions, nil)
		app := fiber.New()
		transaction.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/transactions", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)
		res, _ := io.ReadAll(resp.Body)

		expected := tests.ReadFile(t, "testdata/get-transactions.json")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, string(expected), string(res))
	})

	suite.T().Run("Should ...", func(t *testing.T) {
		suite.service.EXPECT().Get(gomock.Any(), suite.affiliateID).Return(nil, errors.New("test err"))
		app := fiber.New()
		transaction.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/transactions", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)
		res, _ := io.ReadAll(resp.Body)

		expected := tests.ReadFile(t, "testdata/get-transactions-error.json")
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, string(expected), string(res))
	})
}
