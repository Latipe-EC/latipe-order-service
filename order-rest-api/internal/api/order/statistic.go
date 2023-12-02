package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/common/errors"
	dto "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/responses"
	"strings"
)

func (o orderApiHandler) AdminCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := o.orderUsecase.AdminCountingOrderAmount(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) UserCountingOrder(ctx *fiber.Ctx) error {

	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	req.OwnerID = userId

	result, err := o.orderUsecase.UserCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)

}

func (o orderApiHandler) StoreCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeId := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.OwnerID = storeId

	result, err := o.orderUsecase.StoreCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)

}

func (o orderApiHandler) DeliveryCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	deli := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	req.OwnerID = deli

	result, err := o.orderUsecase.DeliveryCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}
