package batch

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type batchResponse struct {
	BatchID uuid.UUID `json:"batchID"`
}

func renderAccepted(c *fiber.Ctx, id uuid.UUID) error {
	return c.Status(http.StatusAccepted).JSON(batchResponse{
		BatchID: id,
	})
}

type fileErrorPresentation struct {
	Row    int      `json:"row"`
	Errors []string `json:"errors,omitempty"`
}
type getBatchesPresentation struct {
	ID          uuid.UUID               `json:"id"`
	AffiliateID uuid.UUID               `json:"affiliateID"`
	Status      string                  `json:"status"`
	Errors      []fileErrorPresentation `json:"errors"`
	CreatedAt   time.Time               `json:"createdAt"`
}

func renderBatches(c *fiber.Ctx, files *[]Model) error {
	if len(*files) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	var items []getBatchesPresentation
	for _, f := range *files {
		item := getBatchesPresentation{
			ID:          f.ID,
			AffiliateID: f.AffiliateID,
			Status:      f.Status,
			CreatedAt:   *f.CreatedAt,
		}

		var fErros []fileErrorPresentation
		_ = json.Unmarshal(*f.Errors, &fErros)
		item.Errors = fErros
		items = append(items, item)
	}
	return c.Status(http.StatusOK).JSON(items)
}
