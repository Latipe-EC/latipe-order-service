package storeserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/storeserv/dto"
)

type Service interface {
	GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error)
}
