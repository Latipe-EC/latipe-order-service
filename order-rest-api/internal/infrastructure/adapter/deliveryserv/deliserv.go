package deliveryserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv/dto"
)

type Service interface {
	CalculateShippingCost(ctx context.Context, req *dto.GetShippingCostRequest) (*dto.GetShippingCostResponse, error)
	GetDeliveryByToken(ctx context.Context, req *dto.GetDeliveryByTokenRequest) (*dto.GetDeliveryByTokenResponse, error)
}
