package batch_test

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"affiliates-backoffice-backend/internal/batch"
	mocks "affiliates-backoffice-backend/internal/batch/mocks"
	"affiliates-backoffice-backend/pkg/log"
	"affiliates-backoffice-backend/pkg/tests"
)

type serviceSuite struct {
	suite.Suite
	ctx         context.Context
	logger      log.LoggerI
	service     batch.ServiceI
	repo        *mocks.MockRepositoryI
	affiliateID string
	files       []batch.Model
}

func (suite *serviceSuite) SetupSuite() {
	loggerTest, _ := log.NewForTest()
	suite.logger = loggerTest
}

func (suite *serviceSuite) BeforeTest(suiteName, testName string) {
	ctrl := gomock.NewController(suite.T())
	suite.repo = mocks.NewMockRepositoryI(ctrl)
	suite.service = batch.NewService(suite.logger, suite.repo)
	suite.ctx = context.Background()

	suite.affiliateID = "ca226967-6b4a-4461-bf16-ab0108f87380"
	dateExpected := "2023-02-26T22:00:16.616853-03:00"
	layout := "2006-01-02T15:04:05-07:00"
	d, err := time.Parse(layout, dateExpected)
	assert.NoError(suite.T(), err)
	suite.files = []batch.Model{
		{
			ID:          uuid.MustParse("4b1362ba-b443-43b4-b057-0d5431a9bbcb"),
			AffiliateID: uuid.MustParse(suite.affiliateID),
			Status:      "CREATED",
			BatchRaw:    "12022-01-15T19:20:30-03:00CURSO DE BEM-ESTAR            0000012750JOSE CARLOS        ",
			Errors:      nil,
			CreatedAt:   &d,
		},
	}
}

func Test_ServiceSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (suite *serviceSuite) Test_Save() {
	suite.T().Run("Should save a batch file", func(t *testing.T) {
		suite.repo.EXPECT().SaveFile(suite.ctx, gomock.Any()).Return(int64(1), nil)
		_, header := makeFormFile(suite.T())

		id, err := suite.service.Save(suite.ctx, "ca226967-6b4a-4461-bf16-ab0108f87380", header)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`), id.String())
	})

	suite.T().Run("Should error if repository fails", func(t *testing.T) {
		affiliateID := "ca226967-6b4a-4461-bf16-ab0108f87380"
		suite.repo.EXPECT().SaveFile(suite.ctx, gomock.Any()).Return(int64(0), errors.New("test error"))
		_, header := makeFormFile(suite.T())

		id, err := suite.service.Save(suite.ctx, affiliateID, header)
		assert.Nil(t, id)
		assert.EqualError(t, err, "test error")
	})
}

func (suite *serviceSuite) Test_Get() {
	suite.T().Run("Should get a batch file", func(t *testing.T) {
		affiliateID := "ca226967-6b4a-4461-bf16-ab0108f87380"
		suite.repo.EXPECT().GetFiles(suite.ctx, affiliateID).Return(&suite.files, nil)

		files, err := suite.service.GetFiles(suite.ctx, affiliateID)
		assert.NoError(t, err)
		assert.Equal(t, &suite.files, files)
	})
}

func makeFormFile(t *testing.T) (multipart.File, *multipart.FileHeader) {
	body, contentType := tests.BuildAttachmentFile(t, true, "testdata/sales.txt", "text/plain")
	req := httptest.NewRequest("POST", "/v1/affiliates/ca226967-6b4a-4461-bf16-ab0108f87380/batches", body)
	req.Header.Add("Content-Type", contentType)
	file, header, err := req.FormFile("attachment")
	assert.NoError(t, err)
	return file, header
}
