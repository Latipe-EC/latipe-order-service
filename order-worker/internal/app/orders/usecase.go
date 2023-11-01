package orders

import (
	"context"
	"order-worker/internal/domain/dto/order"
)

type Usecase interface {
	RollBackQuantity(ctx context.Context, dto *order.OrderCacheData) error
	CreateOrder(ctx context.Context, orderCacheKey string) error
}
