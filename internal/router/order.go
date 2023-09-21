package router

import (
	"github.com/gofiber/fiber/v2"
	"order-service-rest-api/internal/api/order"
)

type OrderRouter interface {
	Init(root *fiber.Router)
}

type orderRouter struct {
	handler order.OrderApiHandler
}

func NewOrderRouter(handler order.OrderApiHandler) OrderRouter {
	return orderRouter{handler: handler}
}

func (o orderRouter) Init(root *fiber.Router) {
	orderRouter := (*root).Group("/orders")
	{
		orderRouter.Post("", o.handler.CreateOrder)
	}
}
