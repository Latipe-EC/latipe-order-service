package deliveryserv

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv/dto"
	http "order-rest-api/pkg/internal_http"
)

var Set = wire.NewSet(
	NewDeliServHttpAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewDeliServHttpAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().SetBaseURL(config.AdapterService.DeliveryService.BaseURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) CalculateShippingCost(ctx context.Context, req *dto.GetShippingCostRequest) (*dto.GetShippingCostResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Shipping Cost]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Shipping Cost]: %s", resp.Body())
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	var regResp *dto.GetShippingCostResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Shipping Cost]: %s", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
