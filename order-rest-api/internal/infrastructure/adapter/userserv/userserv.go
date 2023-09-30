package userserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/userserv/dto"
)

type Service interface {
	Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error)
}
