package userserv

import (
	"context"
	"encoding/json"
	"errors"
	"order-worker/config"
	"order-worker/internal/infrastructure/adapter/userserv/dto"
	http "order-worker/pkg/internal_http"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewUserServHttpAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewUserServHttpAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().SetBaseURL(config.AdapterService.UserService.AuthURL))

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
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", err)
		return nil, errors.New("ErrInternalServer")
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", resp.Body())
		return nil, errors.New("ErrBadRequest")
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", resp.Body())
		return nil, errors.New("ErrInternalServer")
	}

	var regResp *dto.AuthorizationResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", err)
		return nil, err
	}

	return regResp, nil
}
