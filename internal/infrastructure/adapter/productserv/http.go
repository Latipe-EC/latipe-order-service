package productserv

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/common/errors"
	"order-service-rest-api/internal/infrastructure/adapter/productserv/dto"

	http "order-service-rest-api/pkg/internal_http"
	mapper "order-service-rest-api/pkg/util/mapper"
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
			SetHeader("X-INTERNAL-SERVICE", config.AdapterService.ProductService.InternalKey))
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetProductOrderInfo(ctx context.Context, req *dto.OrderProductRequest) (*dto.OrderProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Get(req.URL())

	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Get product]: %s", "ERROR", resp.Body())
		return nil, err
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.OrderProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) ReduceProductQuantity(ctx context.Context, req *dto.ReduceProductRequest) (*dto.ReduceProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Patch(req.URL())

	if err != nil {
		log.Errorf("[Reduce Quantity]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Reduce Quantity]: %s", resp.Body())
		return nil, err
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[Reduce Quantity]: %s", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.ReduceProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) RollBackQuantityOrder(ctx context.Context, req *dto.RollbackQuantityRequest) (*dto.RollbackQuantityResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Patch(req.URL())

	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", resp.Body())
		return nil, err
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.RollbackQuantityResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
