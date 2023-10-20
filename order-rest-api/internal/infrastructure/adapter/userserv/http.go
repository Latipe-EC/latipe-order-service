package userserv

import (
	"context"
	"encoding/json"
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
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", resp.Body())
		return nil, err
	}

	var regResp *dto.AuthorizationResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}
	/*var rawResp dto.BaseResponse*/
	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	/*var regResp *dto.AuthorizationResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Authorize token]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}*/

	return regResp, nil
}
