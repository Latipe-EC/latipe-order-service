package userserv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/userserv/dto"
	http "order-rest-api/pkg/internal_http"
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
			Resty().SetBaseURL(config.AdapterService.UserService.UserURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetAddressDetails(ctx context.Context, req *dto.GetDetailAddressRequest) (*dto.GetDetailAddressResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		Get(req.URL() + fmt.Sprintf("/%v", req.AddressId))

	if err != nil {
		log.Errorf("[Get details address]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Authorize token]: %s", resp.Body())
		return nil, err
	}

	var regResp *dto.GetDetailAddressResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
