package storeserv

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/storeserv/dto"
	"order-rest-api/pkg/internal_http"
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
		return nil, errors.ErrInternalServer
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get store]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	regResp := dto.GetStoreIdByUserResponse{
		StoreID: resp.String(),
	}

	return &regResp, nil
}
