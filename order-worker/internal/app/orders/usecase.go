package orders

import (
	"context"
	"order-worker/internal/domain/dto/order"
)

type Usecase interface {
	RollBackQuantity(ctx context.Context, dto *order.CreateOrderRequest) error
	CreateOrder(ctx context.Context, dto *order.CreateOrderRequest) error
}
