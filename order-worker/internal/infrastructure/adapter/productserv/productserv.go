package productserv

import (
	"order-worker/internal/infrastructure/adapter/productserv/dto"
)

type Service interface {
	GetProductOrderInfo(req *dto.OrderProductRequest) (*dto.OrderProductResponse, error)
	ReduceProductQuantity(req *dto.ReduceProductRequest) (*dto.ReduceProductResponse, error)
	RollBackQuantityOrder(req *dto.RollbackQuantityRequest) (*dto.RollbackQuantityResponse, error)
}
