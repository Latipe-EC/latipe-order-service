package voucherserv

import (
	"context"
	"order-worker/internal/infrastructure/adapter/vouchersev/dto"
)

type Service interface {
	ApplyVoucher(ctx context.Context, req *dto.ApplyVoucherRequest) (*dto.UseVoucherResponse, error)
	Rollback(ctx context.Context, req *dto.RollbackVoucherRequest) (*dto.UseVoucherResponse, error)
}
