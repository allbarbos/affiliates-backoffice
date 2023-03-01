package batch_test

import (
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

	"affiliates-backoffice-backend/internal/batch"
	mocks "affiliates-backoffice-backend/internal/batch/mocks"
	"affiliates-backoffice-backend/pkg/log"
	"affiliates-backoffice-backend/pkg/tests"
)

type handlerSuite struct {
	suite.Suite
	ctx         *fiber.Ctx
	service     *mocks.MockServiceI
	logger      log.LoggerI
	affiliateID string
	files       []batch.Model
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
	dateExpected := "2023-02-26T22:00:16.616853-03:00"
	layout := "2006-01-02T15:04:05-07:00"
	d, err := time.Parse(layout, dateExpected)
	assert.NoError(suite.T(), err)
	errBytes := []byte("[]")
	suite.files = []batch.Model{
		{
			ID:          uuid.MustParse("4b1362ba-b443-43b4-b057-0d5431a9bbcb"),
			AffiliateID: uuid.MustParse(suite.affiliateID),
			Status:      "CREATED",
			BatchRaw:    "12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS        ",
			Errors:      &errBytes,
			CreatedAt:   &d,
		},
	}

	ctrl := gomock.NewController(suite.T())
	suite.service = mocks.NewMockServiceI(ctrl)
}

func Test_HandlerSuite(t *testing.T) {
	suite.Run(t, new(handlerSuite))
}

func (suite *handlerSuite) Test_Get() {
	suite.T().Run("Should get affiliate files", func(t *testing.T) {
		suite.service.EXPECT().GetFiles(gomock.Any(), suite.affiliateID).Return(&suite.files, nil)
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/batches", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)
		res, _ := io.ReadAll(resp.Body)

		expected := tests.ReadFile(t, "testdata/get-affiliate-files.json")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, string(expected), string(res))
	})

	suite.T().Run("Should return NoContent when there is no file", func(t *testing.T) {
		files := []batch.Model{}
		suite.service.EXPECT().GetFiles(gomock.Any(), suite.affiliateID).Return(&files, nil)
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/batches", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	suite.T().Run("Should get affiliate files with errors", func(t *testing.T) {
		*suite.files[0].Errors = tests.ReadFile(t, "testdata/batch-file-errors.json")
		suite.service.EXPECT().GetFiles(gomock.Any(), suite.affiliateID).Return(&suite.files, nil)
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/batches", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)
		res, _ := io.ReadAll(resp.Body)

		expected := tests.ReadFile(t, "testdata/get-affiliate-files-with-err.json")
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, string(expected), string(res))
	})

	suite.T().Run("Should return an error when file fetch fails", func(t *testing.T) {
		suite.service.EXPECT().GetFiles(gomock.Any(), suite.affiliateID).Return(nil, errors.New("test error"))
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/affiliates/%s/batches", suite.affiliateID), nil)
		resp, _ := app.Test(req, -1)
		res, _ := io.ReadAll(resp.Body)

		expected := tests.ReadFile(t, "testdata/get-affiliate-internal-err.json")
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.JSONEq(t, string(expected), string(res))
	})
}

func (suite *handlerSuite) Test_Post() {
	suite.T().Run("Should save a batch file", func(t *testing.T) {
		batchID := uuid.New()
		payloadExpected := fmt.Sprintf(`{"batchID":"%s"}`, batchID)
		suite.service.EXPECT().Save(gomock.Any(), "ca226967-6b4a-4461-bf16-ab0108f87380", gomock.Any()).Return(&batchID, nil)
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		body, contentType := tests.BuildAttachmentFile(t, true, "testdata/sales.txt", "text/plain")
		req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", body)
		req.Header.Add("Content-Type", contentType)

		resp, _ := app.Test(req, -1)
		resBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		assert.JSONEq(t, payloadExpected, string(resBody))
	})

	suite.T().Run("Should error if attachment is not present", func(t *testing.T) {
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)
		tests.BuildAttachmentFile(t, false, "", "")
		req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", nil)

		resp, _ := app.Test(req, -1)
		resBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, `{"errors":["Unable to receive attachment"]}`, string(resBody))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	suite.T().Run("Should error if Content-Type is invalid", func(t *testing.T) {
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		body, contentType := tests.BuildAttachmentFile(t, true, "testdata/sales.txt", "application/json")
		req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", body)
		req.Header.Add("Content-Type", contentType)

		resp, _ := app.Test(req, -1)
		resBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, `{"errors":["Attachment content-type application/json unsupported"]}`, string(resBody))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	suite.T().Run("Should error if attachment content is empty", func(t *testing.T) {
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		body, contentType := tests.BuildAttachmentFile(t, true, "testdata/sales-empty.txt", "text/plain")
		req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", body)
		req.Header.Add("Content-Type", contentType)

		resp, _ := app.Test(req, -1)
		resBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, `{"errors":["Attachment cannot be empty"]}`, string(resBody))
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	suite.T().Run("Should error if save fails", func(t *testing.T) {
		suite.service.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("save fail"))
		app := fiber.New()
		batch.RegisterHandlers(app, suite.logger, suite.service)

		body, contentType := tests.BuildAttachmentFile(t, true, "testdata/sales.txt", "text/plain")
		req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", body)
		req.Header.Add("Content-Type", contentType)

		resp, _ := app.Test(req, -1)
		resBody, _ := io.ReadAll(resp.Body)
		assert.Equal(t, `{"errors":["save fail"]}`, string(resBody))
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
