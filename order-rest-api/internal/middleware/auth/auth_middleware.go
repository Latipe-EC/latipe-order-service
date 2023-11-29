package auth

import (
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/infrastructure/adapter/authserv"
	"order-rest-api/internal/infrastructure/adapter/authserv/dto"
	"order-rest-api/internal/infrastructure/adapter/deliveryserv"
	deliDto "order-rest-api/internal/infrastructure/adapter/deliveryserv/dto"
	"order-rest-api/internal/infrastructure/adapter/storeserv"
	storeDTO "order-rest-api/internal/infrastructure/adapter/storeserv/dto"
	"strings"
)

type AuthenticationMiddleware struct {
	authServ  authserv.Service
	storeServ storeserv.Service
	delivery  deliveryserv.Service
}

func NewAuthMiddleware(authServ authserv.Service, storeServ storeserv.Service, deliServ deliveryserv.Service) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{authServ: authServ, storeServ: storeServ, delivery: deliServ}
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

		//validate store
		storeReq := storeDTO.GetStoreIdByUserRequest{UserID: resp.Id}
		storeReq.BaseHeader.BearToken = bearToken

		storeResp, err := a.storeServ.GetStoreByUserId(ctx.Context(), &storeReq)
		if err != nil {
			return err
		}

		if storeResp.StoreID == "" {
			return errors.ErrPermissionDenied
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		ctx.Locals(STORE_ID, storeResp.StoreID)

		return ctx.Next()
	}
}

func (a AuthenticationMiddleware) RequiredDeliveryAuthentication() fiber.Handler {
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

		//validate store
		deliReq := deliDto.GetDeliveryByTokenRequest{BearerToken: bearToken}

		deliResp, err := a.delivery.GetDeliveryByToken(ctx.Context(), &deliReq)
		if err != nil {
			return err
		}

		if deliResp.ID == "" {
			return errors.ErrPermissionDenied
		}

		ctx.Locals(USERNAME, resp.Email)
		ctx.Locals(USER_ID, resp.Id)
		ctx.Locals(ROLE, resp.Role)
		ctx.Locals(BEARER_TOKEN, bearToken)
		ctx.Locals(DELIVERY_ID, deliResp.ID)

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
