package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"order-service-rest-api/internal/app/orders"
	"order-service-rest-api/internal/common/errors"
	dto "order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/message"
	"order-service-rest-api/internal/middleware/auth"
	"order-service-rest-api/internal/responses"
	"order-service-rest-api/pkg/util/valid"
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

	if err := ctx.BodyParser(&bodyReq); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := valid.GetValidator().Validate(bodyReq); err != nil {
		return errors.ErrBadRequest
	}

	bodyReq.UserRequest.UserId = userId
	bodyReq.UserRequest.Username = username

	if err := o.orderUsecase.PendingOrder(context, &bodyReq); err != nil {
		return err
	}

	if err := message.SendMessage(bodyReq); err != nil {
		resp := responses.DefaultError
		resp.Status = 500
		resp.Message = "Fail"
		return resp.JSON(ctx)
	}

	return responses.DefaultSuccess.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o orderApiHandler) ListOfOrder(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (o orderApiHandler) GetOrderById(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
