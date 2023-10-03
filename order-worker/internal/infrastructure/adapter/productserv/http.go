package productserv

import (
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
			Resty().SetDebug(true).
			SetBaseURL(config.AdapterService.ProductService.BaseURL).
			SetHeader("X-INTERNAL-SERVICE", config.AdapterService.ProductService.InternalKey))
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetProductOrderInfo(req *dto.OrderProductRequest) (*dto.OrderProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req.Items).
		Post(req.URL())

	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Get product]: %s", "ERROR", resp.Body())
		return nil, errors.New("internal service request")
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[%s] [Get product]: %s", "ERROR", resp.Body())
		return nil, errors.New("service bad request")
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp.Data); err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, err
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.OrderProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Get product]: %s", "ERROR", err)
		return nil, err
	}

	return regResp, nil
}

func (h httpAdapter) ReduceProductQuantity(req *dto.ReduceProductRequest) (*dto.ReduceProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
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
		return nil, err
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.ReduceProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, err
	}

	return regResp, nil
}

func (h httpAdapter) RollBackQuantityOrder(req *dto.RollbackQuantityRequest) (*dto.RollbackQuantityResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
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
		return nil, err
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *dto.RollbackQuantityResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, err
	}

	return regResp, nil
}
