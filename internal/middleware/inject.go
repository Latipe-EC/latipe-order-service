package middleware

import (
	"github.com/google/wire"
	"order-service-rest-api/internal/middleware/auth"
)

var Set = wire.NewSet(
	NewMiddleware,
	auth.NewAuthMiddleware,
)
