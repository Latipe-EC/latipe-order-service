package order

import (
	"order-rest-api/pkg/util/pagable"
)

type Repository interface {
	FindById(Id int) (*Order, error)
	FindByUUID(uuid string) (*Order, error)
	FindOrderByStoreID(storeId string, query *pagable.Query) ([]Order, error)
	FindAll(query *pagable.Query) ([]Order, error)
	FindByUserId(userId string, query *pagable.Query) ([]Order, error)
	FindOrderLogByOrderId(orderId int) ([]OrderStatusLog, error)
	FindOrderByUserAndProduct(userId string, productId string) ([]Order, error)
	Save(order *Order) error
	Update(order Order) error
	UpdateStatus(orderID int, status int) error
	Total(query *pagable.Query) (int, error)
}
