package storeserv

import (
	"context"
	"order-worker/internal/infrastructure/adapter/storeserv/dto"
)

type Service interface {
	GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error)
	GetStoreByStoreId(ctx context.Context, req *dto.GetStoreByIdRequest) (*dto.GetStoreByIdResponse, error)
}
