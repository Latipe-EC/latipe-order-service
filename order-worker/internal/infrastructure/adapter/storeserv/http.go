package storeserv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-worker/config"
	"order-worker/internal/infrastructure/adapter/storeserv/dto"
	http "order-worker/pkg/internal_http"
	"order-worker/pkg/util/mapper"
)

var Set = wire.NewSet(
	NewStoreServiceAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewStoreServiceAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().
			SetBaseURL(config.AdapterService.StoreService.BaseURL).
			SetHeader("X-INTERNAL-SERVICE", config.AdapterService.StoreService.InternalKey))
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BaseHeader.BearToken)).
		Get(req.URL() + req.UserID)

	if err != nil {
		log.Errorf("[Get store]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Get store]: %s", resp.Body())
		return nil, errors.New("internal service request")
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get store]: %s", resp.Body())
		return nil, errors.New("service bad request")
	}

	regResp := dto.GetStoreIdByUserResponse{
		StoreID: resp.String(),
	}

	return &regResp, nil
}

func (h httpAdapter) GetStoreByStoreId(ctx context.Context, req *dto.GetStoreByIdRequest) (*dto.GetStoreByIdResponse, error) {
	resp, err := h.client.MakeRequest().
		Get(fmt.Sprintf("%v%v", req.URL(), req.StoreID))

	if err != nil {
		log.Errorf("[Get store id]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Get store id]: %s", resp.Body())
		return nil, errors.New("internal service request")
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get store id]: %s", resp.Body())
		return nil, errors.New("service bad request")
	}

	var rawResp dto.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp.Data); err != nil {
		log.Errorf("[%s] [Get store id]: %s", "ERROR", err)
		return nil, err
	}

	var regResp *dto.GetStoreByIdResponse

	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[Get store id]: %s", err)
		return nil, err
	}

	return regResp, nil
}
