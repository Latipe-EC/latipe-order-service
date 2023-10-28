package userserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/userserv/dto"
)

type Service interface {
	GetAddressDetails(ctx context.Context, req *dto.GetDetailAddressRequest) (*dto.GetDetailAddressResponse, error)
}
