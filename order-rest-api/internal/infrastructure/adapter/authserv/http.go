package authserv

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/authserv/dto"
	http "order-rest-api/pkg/internal_http"
)

var Set = wire.NewSet(
	NewAuthServHttpAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewAuthServHttpAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().SetBaseURL(config.AdapterService.AuthService.BaseURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Authorize token]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Authorize token]: %s", "ERROR", resp.Body())
		return nil, err
	}

	var regResp *dto.AuthorizationResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Authorize token]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
