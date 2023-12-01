package order

import (
	"order-worker/internal/domain/entities/custom"
	"order-worker/pkg/util/pagable"
)

type Repository interface {
	GetOrderAmountOfStore(orderId int) ([]custom.AmountItemOfStoreInOrder, error)
	FindAllFinishShippingOrder() ([]Order, error)
	CreateOrderCommmsionTransaction(order *Order, ocms *OrderCommission, log *OrderStatusLog) error
	Save(order *Order) error
	FindById(Id string) (*Order, error)
	FindByUserId(query *pagable.Query) ([]*Order, error)
	UpdateOrderRating(itemId string, ratingId string) error
	Update(order Order) error
	Total(query *pagable.Query) (int, error)
}
