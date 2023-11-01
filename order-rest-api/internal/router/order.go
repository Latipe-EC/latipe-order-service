package router

import (
	"github.com/gofiber/fiber/v2"
	"order-rest-api/internal/api/order"
	"order-rest-api/internal/middleware"
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
		orderRouter.Get("/:id", o.middleware.Authentication.RequiredAuthentication(), o.handler.GetOrderById)
		orderRouter.Get("", o.middleware.Authentication.RequiredAuthentication(), o.handler.ListOfOrder)
		orderRouter.Patch("/:id/status", o.middleware.Authentication.RequiredAuthentication(), o.handler.CreateOrder)
	}
	products := orderRouter.Group("/products")
	{
		products.Post("/purchased", o.middleware.Authentication.RequiredAuthentication(), o.handler.CheckOrderOfUser)
	}

}
