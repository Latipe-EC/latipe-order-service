package orders

import (
	"context"
	"order-worker/internal/domain/dto/order"
)

type Usecase interface {
	CreateOrder(ctx context.Context, data *order.OrderMessage) error
	CreateCommissionOrderComplete(ctx context.Context) error
}
