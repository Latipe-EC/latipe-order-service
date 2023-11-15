package auth

import (
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/authserv"
	"order-rest-api/internal/infrastructure/adapter/authserv/dto"
	"strings"
)

type AuthenticationMiddleware struct {
	authServ authserv.Service
}

func NewAuthMiddleware(authServ authserv.Service) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{authServ: authServ}
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
		resp, err := a.authServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return err
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredStoreAuthentication() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}
		resp, err := a.authServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return err
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredRole(roles []string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		bearToken := ctx.Get("Authorization")
		if bearToken == "" {
			return errors.ErrUnauthenticated
		}

		bearToken = strings.Split(bearToken, " ")[1]
		req := dto.AuthorizationRequest{
			Token: bearToken,
		}
		resp, err := a.authServ.Authorization(ctx.Context(), &req)
		if err != nil {
			return err
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		return ctx.Next()
	}
}
