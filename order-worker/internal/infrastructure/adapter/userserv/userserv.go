package userserv

import (
	"context"
	"order-worker/internal/infrastructure/adapter/userserv/dto"
)

type Service interface {
	Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error)
}
