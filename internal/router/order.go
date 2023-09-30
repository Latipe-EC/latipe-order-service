package router

import (
	"github.com/gofiber/fiber/v2"
	"order-service-rest-api/internal/api/order"
	"order-service-rest-api/internal/middleware"
)

type OrderRouter interface {
	Init(root *fiber.Router)
}

type orderRouter struct {
	handler    order.OrderApiHandler
	middleware *middleware.Middleware
}

func NewOrderRouter(handler order.OrderApiHandler, middleware *middleware.Middleware) OrderRouter {
	return orderRouter{
		handler:    handler,
		middleware: middleware,
	}
}

func (o orderRouter) Init(root *fiber.Router) {
	orderRouter := (*root).Group("/orders")
	{
		orderRouter.Post("", o.middleware.Authentication.RequiredAuthentication(), o.handler.CreateOrder)
	}
}
