package productserv

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	productDTO "order-rest-api/internal/infrastructure/adapter/productserv/dto"
	"order-rest-api/pkg/internal_http"
	"order-rest-api/pkg/util/mapper"
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

func (h httpAdapter) GetProductOrderInfo(ctx context.Context, req *productDTO.OrderProductRequest) (*productDTO.OrderProductResponse, error) {
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

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *productDTO.OrderProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) ReduceProductQuantity(ctx context.Context, req *productDTO.ReduceProductRequest) (*productDTO.ReduceProductResponse, error) {
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

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[Reduce Quantity]: %s", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *productDTO.ReduceProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) RollBackQuantityOrder(ctx context.Context, req *productDTO.RollbackQuantityRequest) (*productDTO.RollbackQuantityResponse, error) {
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

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *productDTO.RollbackQuantityResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
