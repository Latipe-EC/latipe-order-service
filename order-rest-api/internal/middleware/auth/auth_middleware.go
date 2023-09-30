package auth

import (
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/userserv"
	"order-rest-api/internal/infrastructure/adapter/userserv/dto"
	"strings"
)

type AuthenticationMiddleware struct {
	userServ userserv.Service
}

func NewAuthMiddleware(userServ userserv.Service) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{userServ: userServ}
}

func (a AuthenticationMiddleware) RequiredAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}
		resp, err := a.userServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return err
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		return ctx.Next()
	}
}
