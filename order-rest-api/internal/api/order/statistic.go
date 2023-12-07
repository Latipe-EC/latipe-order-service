package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/app/orders"
	"order-rest-api/internal/common/errors"
	"order-rest-api/internal/domain/dto/order/statistic"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/responses"
	"order-rest-api/pkg/util/valid"
	"strings"
)

type statisticApiHandler struct {
	orderUsecase orders.Usecase
}

func NewStatisticHandler(orderUsecase orders.Usecase) OrderStatisticApiHandler {
	return statisticApiHandler{
		orderUsecase: orderUsecase,
	}
}

func (s statisticApiHandler) AdminGetTotalOrderInSystemInDay(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminTotalOrderInDayRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderUsecase.AdminGetTotalOrderInSystemInDay(context, &req)
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

func (s statisticApiHandler) AdminGetTotalOrderInSystemInMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminTotalOrderInMonthRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderUsecase.AdminGetTotalOrderInSystemInMonth(context, &req)
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

func (s statisticApiHandler) AdminGetTotalOrderInSystemInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.AdminGetTotalOrderInYearRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := s.orderUsecase.AdminGetTotalOrderInSystemInYear(context, &req)
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

func (s statisticApiHandler) AdminGetTotalCommissionOrderInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.OrderCommissionDetailRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderUsecase.AdminGetTotalCommissionOrderInYear(context, &req)
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

func (s statisticApiHandler) AdminListOfProductSoldOnMonth(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ListOfProductSoldRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	result, err := s.orderUsecase.AdminListOfProductSoldOnMonth(context, &req)
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

func (s statisticApiHandler) GetTotalOrderInMonthOfStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetTotalStoreOrderInMonthRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderUsecase.GetTotalOrderInMonthOfStore(context, &req)
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

func (s statisticApiHandler) GetTotalOrderInYearOfStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.GetTotalOrderInYearOfStoreRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := s.orderUsecase.GetTotalOrderInYearOfStore(context, &req)
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

func (s statisticApiHandler) GetTotalStoreCommissionInYear(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.OrderCommissionDetailRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderUsecase.GetTotalStoreCommissionInYear(context, &req)
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

func (s statisticApiHandler) ListOfProductSoldOnMonthStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := statistic.ListOfProductSoldRequest{}

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Date == "" {
		req.Date = InitDateValue()
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreId = storeID

	result, err := s.orderUsecase.ListOfProductSoldOnMonthStore(context, &req)
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
