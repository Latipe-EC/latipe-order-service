package orders

import (
	"context"
	"order-service-rest-api/internal/domain/dto/order"
)

type Usecase interface {
	PendingOrder(ctx context.Context, dto *order.CreateOrderRequest) error
	RollBackQuantity(ctx context.Context, dto *order.CreateOrderRequest) error
	CreateOrder(ctx context.Context, dto *order.CreateOrderRequest) error
	GetOrderById(ctx context.Context, dto *order.GetOrderRequest) (*order.GetOrderResponse, error)
	GetOrderList(ctx context.Context, dto *order.GetOrderListRequest) (*order.GetOrderListResponse, error)
	UpdateStatusOrder(ctx context.Context, dto *order.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *order.UpdateOrderRequest) error
}
