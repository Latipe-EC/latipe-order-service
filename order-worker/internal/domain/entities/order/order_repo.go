package order

import (
	"order-worker/internal/domain/entities/custom"
	"order-worker/pkg/util/pagable"
)

type Repository interface {
	GetOrderAmountOfStore(orderId int) ([]custom.AmountItemOfStoreInOrder, error)
	FindAllFinishShippingOrder() ([]Order, error)
	PendingOrderShippingOrder() ([]Order, error)

	UpdateCommission(Id string) (*Order, error)
	Save(order *Order) error
	FindByUUID(Id string) (*Order, error)
	FindByID(Id int) (*Order, error)
	FindByUserId(query *pagable.Query) ([]*Order, error)
	UpdateOrderRating(itemId string, ratingId string) error
	Update(order Order) error
	Total(query *pagable.Query) (int, error)

	//cms
	UpdateOrderCommissionTransaction(order *Order, ocms *OrderCommission, log *OrderStatusLog) error
	CreateOrderCommissionTransaction(ocms *OrderCommission) error
	FindCommissionByOrderId(orderId int) (*OrderCommission, error)
}
