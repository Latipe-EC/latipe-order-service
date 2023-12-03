package voucherserv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-rest-api/config"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/vouchersev/dto"

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
			Resty().SetBaseURL(config.AdapterService.PromotionService.BaseURL))

	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		SetBody(req).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Apply voucher]: %s", err)
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Apply voucher]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Apply voucher]: %s", resp.Body())
		return nil, errors.ErrInternalServer
	}

	var regResp *dto.UseVoucherResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[%s] [Apply voucher]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
