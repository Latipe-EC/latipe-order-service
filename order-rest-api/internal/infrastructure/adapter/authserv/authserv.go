package authserv

import (
	"context"
	"order-rest-api/internal/infrastructure/adapter/authserv/dto"
)

type Service interface {
	Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error)
}
