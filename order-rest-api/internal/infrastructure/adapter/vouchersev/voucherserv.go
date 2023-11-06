package voucherserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/vouchersev/dto"
)

type Service interface {
	ApplyVoucher(ctx context.Context, req *dto.ApplyVoucherRequest) (*dto.UseVoucherResponse, error)
}
