package orders

import (
	"context"
	"order-worker/internal/domain/dto/order"
)

type Usecase interface {
	RollBackQuantity(ctx context.Context, dto *order.OrderMessage) error
	CreateOrder(ctx context.Context, data *order.OrderMessage) error
}
