package orders

import (
	"context"
	"order-service-rest-api/internal/domain/dto/order"
)

type Usecase interface {
	CreateOrder(ctx context.Context, dto *order.CreateOrderRequest) error
	GetOrderById(ctx context.Context, dto *order.GetOrderRequest) (*order.GetOrderResponse, error)
	GetOrderList(ctx context.Context, dto *order.GetOrderListRequest) (*order.GetOrderListResponse, error)
	UpdateStatusOrder(ctx context.Context, dto *order.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *order.UpdateOrderRequest) error
}
