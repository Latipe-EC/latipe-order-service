package order

import (
	"order-rest-api/pkg/util/pagable"
)

type Repository interface {
	FindById(Id int) (*Order, error)
	FindAll(query *pagable.Query) ([]Order, error)
	FindByUserId(userId int, query *pagable.Query) ([]Order, error)
	FindOrderLogByOrderId(orderId int) ([]OrderStatusLog, error)
	Save(order *Order) error
	Update(order Order) error
	Total(query *pagable.Query) (int, error)
}
