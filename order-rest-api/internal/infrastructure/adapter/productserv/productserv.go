package productserv

import (
	"context"
	dto2 "order-rest-api/internal/infrastructure/adapter/productserv/dto"
)

type Service interface {
	GetProductOrderInfo(ctx context.Context, req *dto2.OrderProductRequest) (*dto2.OrderProductResponse, error)
	ReduceProductQuantity(ctx context.Context, req *dto2.ReduceProductRequest) (*dto2.ReduceProductResponse, error)
	RollBackQuantityOrder(ctx context.Context, req *dto2.RollbackQuantityRequest) (*dto2.RollbackQuantityResponse, error)
}
