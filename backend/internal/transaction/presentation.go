package transaction

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type transactions struct {
	BatchID     uuid.UUID `json:"batchID"`
	AffiliateID uuid.UUID `json:"affiliateID"`
	Type        int       `json:"type"`
	Date        time.Time `json:"date"`
	Product     string    `json:"product"`
	Value       float64   `json:"value"`
	Seller      string    `json:"seller"`
}

func renderTransactions(c *fiber.Ctx, trans []Model) error {
	t := make([]transactions, len(trans))
	for i, tran := range trans {
		t[i].BatchID = tran.BatchID
		t[i].AffiliateID = tran.AffiliateID
		t[i].Type = tran.Type
		t[i].Date = tran.Date
		t[i].Product = tran.Product
		t[i].Value = tran.Value
		t[i].Seller = tran.Seller
	}
	return c.
		Status(http.StatusOK).
		JSON(t)
}
