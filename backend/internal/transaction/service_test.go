package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"affiliates-backoffice-backend/internal/transaction"
	mocks "affiliates-backoffice-backend/internal/transaction/mocks"
	"affiliates-backoffice-backend/pkg/log"
)

type serviceSuite struct {
	suite.Suite
	ctx         context.Context
	logger      log.LoggerI
	service     transaction.ServiceI
	repo        *mocks.MockRepositoryI
	affiliateID string
}

func (suite *serviceSuite) SetupSuite() {
	loggerTest, _ := log.NewForTest()
	suite.logger = loggerTest
}

func (suite *serviceSuite) BeforeTest(suiteName, testName string) {
	ctrl := gomock.NewController(suite.T())
	suite.repo = mocks.NewMockRepositoryI(ctrl)
	suite.service = transaction.NewService(suite.logger, suite.repo)
	suite.ctx = context.Background()
	suite.affiliateID = "ca226967-6b4a-4461-bf16-ab0108f87380"
}

func Test_ServiceSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (ss *serviceSuite) Test_Get() {
	ss.T().Run("Should get a transactions", func(t *testing.T) {
		items := []transaction.Model{
			{
				BatchID:     uuid.New(),
				AffiliateID: uuid.MustParse(ss.affiliateID),
				Type:        1,
				Date:        time.Now(),
				Product:     "Product",
				Value:       1.40,
				Seller:      "Seller",
			},
		}
		ss.repo.EXPECT().GetTransactions(ss.affiliateID).Return(&items, nil)

		actual, err := ss.service.Get(context.Background(), ss.affiliateID)
		a := *actual
		assert.NoError(t, err)
		assert.Equal(t, items[0], a[0])
	})
}
