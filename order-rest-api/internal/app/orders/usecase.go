package orders

import (
	"context"
	order2 "order-rest-api/internal/domain/dto/order"
)

type Usecase interface {
	PendingOrder(ctx context.Context, dto *order2.CreateOrderRequest) error
	RollBackQuantity(ctx context.Context, dto *order2.CreateOrderRequest) error
	CreateOrder(ctx context.Context, dto *order2.CreateOrderRequest) error
	GetOrderById(ctx context.Context, dto *order2.GetOrderRequest) (*order2.GetOrderResponse, error)
	GetOrderList(ctx context.Context, dto *order2.GetOrderListRequest) (*order2.GetOrderListResponse, error)
	UpdateStatusOrder(ctx context.Context, dto *order2.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *order2.UpdateOrderRequest) error
}
