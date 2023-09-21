package order

import (
	"github.com/gofiber/fiber/v2"
	"order-service-rest-api/internal/common/errors"
	dto "order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/domain/entities/order"
	"order-service-rest-api/internal/message"
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
	orderRepos order.Repository
}

func NewOrderHandler(orderRepo order.Repository) OrderApiHandler {
	return orderApiHandler{
		orderRepos: orderRepo,
	}
}

func (o orderApiHandler) CreateOrder(ctx *fiber.Ctx) error {
	bodyReq := dto.CreateOrderRequest{}

	if err := ctx.BodyParser(&bodyReq); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := valid.GetValidator().Validate(bodyReq); err != nil {
		return errors.ErrBadRequest
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
