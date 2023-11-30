package productserv

import (
	"context"
	"order-worker/internal/infrastructure/adapter/productserv/dto"
)

type Service interface {
	GetProductOrderInfo(ctx context.Context, req *dto.OrderProductRequest) (*dto.OrderProductResponse, error)
	UpdateProductQuantity(ctx context.Context, req *dto.ReduceProductRequest) error
}
