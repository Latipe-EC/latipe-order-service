package middleware

import (
	"github.com/google/wire"
	"order-rest-api/internal/middleware/auth"
)

var Set = wire.NewSet(
	NewMiddleware,
	auth.NewAuthMiddleware,
)
