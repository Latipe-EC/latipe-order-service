package middleware

import "order-service-rest-api/internal/middleware/auth"

type Middleware struct {
	Authentication *auth.AuthenticationMiddleware
}

func NewMiddleware(auth *auth.AuthenticationMiddleware) *Middleware {
	return &Middleware{Authentication: auth}
}
