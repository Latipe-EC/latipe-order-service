package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"order-rest-api/internal/app/orders"
	"order-rest-api/internal/common/errors"
	dto "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/message"
	"order-rest-api/internal/middleware/auth"
	"order-rest-api/internal/responses"
	"order-rest-api/pkg/util/pagable"
	"order-rest-api/pkg/util/valid"
	"strings"
)

type OrderApiHandler interface {
	CreateOrder(ctx *fiber.Ctx) error
	UpdateOrderStatus(ctx *fiber.Ctx) error
	ListOfOrder(ctx *fiber.Ctx) error
	GetOrderById(ctx *fiber.Ctx) error
}

type orderApiHandler struct {
	orderUsecase orders.Usecase
}

func NewOrderHandler(orderUsecase orders.Usecase) OrderApiHandler {
	return orderApiHandler{
		orderUsecase: orderUsecase,
	}
}

func (o orderApiHandler) CreateOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	bodyReq := dto.CreateOrderRequest{}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	username := fmt.Sprintf("%v", ctx.Locals(auth.USERNAME))
	if username == "" {
		return errors.ErrUnauthenticated
	}

	bearerToken := fmt.Sprintf("%v", ctx.Locals(auth.BEARER_TOKEN))
	if bearerToken == "" {
		return errors.ErrUnauthenticated
	}

	if err := ctx.BodyParser(&bodyReq); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := valid.GetValidator().Validate(bodyReq); err != nil {
		return errors.ErrBadRequest
	}

	bodyReq.Header.BearerToken = bearerToken
	bodyReq.UserRequest.UserId = userId
	bodyReq.UserRequest.Username = username

	orderKey, err := o.orderUsecase.ProcessCacheOrder(context, &bodyReq)
	if err != nil {
		return err
	}

	if err := message.SendMessage(orderKey); err != nil {
		resp := responses.DefaultError
		resp.Status = 500
		resp.Message = "Fail"
		return resp.JSON(ctx)
	}

	dataResp := make(map[string]string)
	dataResp["order_key"] = orderKey
	dataResp["message"] = "The order was created successful"

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o orderApiHandler) ListOfOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(dto.GetOrderListRequest)
	req.Query = query

	result, err := o.orderUsecase.GetOrderList(context, req)
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

func (o orderApiHandler) GetOrderById(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderUsecase.GetOrderById(context, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}
