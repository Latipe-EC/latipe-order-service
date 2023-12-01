package orders

import (
	"context"
	"order-worker/internal/domain/dto"
	"order-worker/internal/domain/dto/order"
)

type Usecase interface {
	CreateOrderTransaction(ctx context.Context, data *order.OrderMessage) error
	CreateCommissionOrderComplete(ctx context.Context) error
	UpdateRatingItem(ctx context.Context, data *dto.RatingMessage) error
}
