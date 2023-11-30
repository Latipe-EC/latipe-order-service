package productserv

import (
	"context"
	"encoding/json"
	"errors"
	"order-worker/config"
	"order-worker/internal/infrastructure/adapter/productserv/dto"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"

	http "order-worker/pkg/internal_http"
	mapper "order-worker/pkg/util/mapper"
)

var Set = wire.NewSet(
	NewProductServAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewProductServAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().
			SetBaseURL(config.AdapterService.ProductService.BaseURL).
			SetHeader("X-API-KEY", config.AdapterService.ProductService.InternalKey))
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetProductOrderInfo(ctx context.Context, req *dto.OrderProductRequest) (*dto.OrderProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetBody(req.Items).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Get product]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.New("internal service request")
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.New("service bad request")
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp.Data); err != nil {
		log.Errorf(" [Get product]: %s", err)
		return nil, err
	}

	var regResp *dto.OrderProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[Get product]: %s", err)
		return nil, err
	}

	return regResp, nil
}

func (h httpAdapter) UpdateProductQuantity(ctx context.Context, req *dto.ReduceProductRequest) error {
	resp, err := h.client.MakeRequest().
		SetBody(req.Items).
		SetContext(ctx).
		Patch(req.URL())

	if err != nil {
		log.Errorf("[Update Quantity]: %s", err)
		return err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Update Quantity: %s", resp.Body())
		return errors.New("service bad request")
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Update Quantity]: %s", resp.Body())
		return err
	}

	log.Infof("Update quantity:%v", req.Items)

	return nil
}
