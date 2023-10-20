package orders

import (
	"context"
	orderDTO "order-rest-api/internal/domain/dto/order"
)

type Usecase interface {
	ProcessCacheOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (string, error)
	GetOrderById(ctx context.Context, dto *orderDTO.GetOrderRequest) (*orderDTO.GetOrderResponse, error)
	GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error)
	GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error)
	UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error
}
