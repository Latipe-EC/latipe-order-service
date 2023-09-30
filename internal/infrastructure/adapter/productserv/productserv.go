package productserv

import (
	"context"
	"order-service-rest-api/internal/infrastructure/adapter/productserv/dto"
)

type Service interface {
	GetProductOrderInfo(ctx context.Context, req *dto.OrderProductRequest) (*dto.OrderProductResponse, error)
	ReduceProductQuantity(ctx context.Context, req *dto.ReduceProductRequest) (*dto.ReduceProductResponse, error)
	RollBackQuantityOrder(ctx context.Context, req *dto.RollbackQuantityRequest) (*dto.RollbackQuantityResponse, error)
}
