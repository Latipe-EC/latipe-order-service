package voucherserv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"order-worker/config"
	"order-worker/internal/infrastructure/adapter/vouchersev/dto"
	http "order-worker/pkg/internal_http"
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

func (h httpAdapter) ApplyVoucher(ctx context.Context, req *dto.ApplyVoucherRequest) (*dto.UseVoucherResponse, error) {
	resp, err := h.client.MakeRequest().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		SetContext(ctx).
		SetBody(req).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Apply voucher]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Apply voucher]: %s", resp.Body())
		return nil, errors.New("bad request")
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Apply voucher]: %s", resp.Body())
		return nil, errors.New("internal server")
	}

	var regResp *dto.UseVoucherResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Apply voucher]: %s", err)
		return nil, err
	}

	log.Infof("[Apply voucher]: %v", req.Vouchers)
	return regResp, nil
}

func (h httpAdapter) Rollback(ctx context.Context, req *dto.RollbackVoucherRequest) (*dto.UseVoucherResponse, error) {
	resp, err := h.client.MakeRequest().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", req.BearerToken)).
		SetBody(req).
		Post(req.URL())

	if err != nil {
		log.Errorf("[rollback voucher]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[rollback voucher]: %s", resp.Body())
		return nil, errors.New("bad request")
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[rollback voucher]: %s", resp.Body())
		return nil, errors.New("internal server")
	}

	var regResp *dto.UseVoucherResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[rollback voucher]: %s", err)
		return nil, err
	}

	return regResp, nil
}
