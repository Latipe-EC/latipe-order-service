package order

import (
	"order-rest-api/internal/domain/dto/custom_entity"
	"order-rest-api/pkg/util/pagable"
)

type Repository interface {
	FindById(Id int) (*Order, error)
	FindByUUID(uuid string) (*Order, error)
	FindOrderByStoreID(storeId string, query *pagable.Query) ([]Order, error)
	FindOrderByDelivery(deliID string, query *pagable.Query) ([]Order, error)
	FindAll(query *pagable.Query) ([]Order, error)
	FindByUserId(userId string, query *pagable.Query) ([]Order, error)
	FindOrderLogByOrderId(orderId int) ([]OrderStatusLog, error)
	FindOrderByUserAndProduct(userId string, productId string) ([]Order, error)
	GetOrderAmountOfStore(orderId int) ([]custom_entity.AmountItemOfStoreInOrder, error)
	Save(order *Order) error
	Update(order Order) error
	UpdateStatus(orderID int, status int) error
	UpdateOrderItem(orderItem string, status int) error
	Total(query *pagable.Query) (int, error)

	//custom_entity - admin
	GetTotalOrderInSystemInDay(date string) ([]custom_entity.TotalOrderInSystemInHours, error)
	GetTotalOrderInSystemInMonth(month int, year int) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYear(year int) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYear(month int, year int, count int) ([]custom_entity.OrderCommissionDetail, error)
	ListOfProductSelledOnMonth(month int, year int, count int) ([]custom_entity.TopOfProductSold, error)

	//custom_entity - store
	GetTotalOrderInSystemInMonthOfStore(month int, year int, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYearOfStore(year int, storeId string) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYearOfStore(month int, year int, count int, storeId string) ([]custom_entity.OrderCommissionDetail, error)
	ListOfProductSelledOnMonthStore(month int, year int, count int, storeId string) ([]custom_entity.TopOfProductSold, error)
}
