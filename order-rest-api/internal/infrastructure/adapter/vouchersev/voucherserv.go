package voucherserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/vouchersev/dto"
)

type Service interface {
	CheckingVoucher(ctx context.Context, req *dto.CheckingVoucherRequest) (*dto.UseVoucherResponse, error)
}
