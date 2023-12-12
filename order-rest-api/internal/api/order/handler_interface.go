package order

import (
	"github.com/gofiber/fiber/v2"
)

type OrderApiHandler interface {
	CreateOrder(ctx *fiber.Ctx) error
	UpdateOrderStatus(ctx *fiber.Ctx) error
	ListOfOrder(ctx *fiber.Ctx) error
	GetOrderByUUID(ctx *fiber.Ctx) error
	GetMyOrder(ctx *fiber.Ctx) error
	UserCancelOrder(ctx *fiber.Ctx) error
	UserRefundOrder(ctx *fiber.Ctx) error
	AdminCancelOrder(ctx *fiber.Ctx) error

	GetMyStoreOrder(ctx *fiber.Ctx) error
	GetStoreOrderDetail(ctx *fiber.Ctx) error
	UpdateOrderItemStatus(ctx *fiber.Ctx) error
	UpdateStatusByDelivery(ctx *fiber.Ctx) error
	CancelOrderItemStatus(ctx *fiber.Ctx) error

	GetOrdersByDelivery(ctx *fiber.Ctx) error

	InternalGetOrderByUUID(ctx *fiber.Ctx) error

	AdminCountingOrder(ctx *fiber.Ctx) error
	UserCountingOrder(ctx *fiber.Ctx) error
	StoreCountingOrder(ctx *fiber.Ctx) error
	DeliveryCountingOrder(ctx *fiber.Ctx) error
	UserGetOrderByUUID(ctx *fiber.Ctx) error
	DeliveryGetOrderByUUID(ctx *fiber.Ctx) error
	SearchOrderIdByKeyword(ctx *fiber.Ctx) error
}

type OrderStatisticApiHandler interface {
	//admin
	AdminGetTotalOrderInSystemInDay(ctx *fiber.Ctx) error
	AdminGetTotalOrderInSystemInMonth(ctx *fiber.Ctx) error
	AdminGetTotalOrderInSystemInYear(ctx *fiber.Ctx) error
	AdminGetTotalCommissionOrderInYear(ctx *fiber.Ctx) error
	AdminListOfProductSoldOnMonth(ctx *fiber.Ctx) error
	//store
	GetTotalOrderInMonthOfStore(ctx *fiber.Ctx) error
	GetTotalOrderInYearOfStore(ctx *fiber.Ctx) error
	GetTotalStoreCommissionInYear(ctx *fiber.Ctx) error
	ListOfProductSoldOnMonthStore(ctx *fiber.Ctx) error
}
