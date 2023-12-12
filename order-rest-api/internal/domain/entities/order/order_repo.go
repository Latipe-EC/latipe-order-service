package order

import (
	"context"
	"order-rest-api/internal/domain/dto/custom_entity"
	"order-rest-api/pkg/util/pagable"
)

type Repository interface {
	FindById(ctx context.Context, Id int) (*Order, error)
	FindByItemId(ctx context.Context, itemId string) (*OrderItem, error)
	FindByUUID(ctx context.Context, uuid string) (*Order, error)
	FindOrderByStoreID(ctx context.Context, storeId string, query *pagable.Query, keyword string) ([]Order, error)
	FindOrderByDelivery(ctx context.Context, deliID string, keyword string, query *pagable.Query) ([]Order, error)
	FindAll(ctx context.Context, query *pagable.Query) ([]Order, error)
	FindByUserId(ctx context.Context, userId string, query *pagable.Query) ([]Order, error)
	SearchOrderByStoreID(ctx context.Context, storeId string, keyword string, query *pagable.Query) ([]Order, error)
	FindOrderLogByOrderId(ctx context.Context, orderId int) ([]OrderStatusLog, error)
	FindOrderByUserAndProduct(ctx context.Context, userId string, productId string) ([]Order, error)
	GetOrderAmountOfStore(ctx context.Context, orderId int) ([]custom_entity.AmountItemOfStoreInOrder, error)
	Save(ctx context.Context, order *Order) error
	Update(ctx context.Context, order Order) error
	UpdateStatus(ctx context.Context, orderID int, status int, message ...string) error
	UpdateOrderItem(ctx context.Context, orderItem string, status int) error
	Total(ctx context.Context, query *pagable.Query) (int, error)
	TotalStoreOrder(ctx context.Context, storeId string, query *pagable.Query, keyword string) (int, error)
	TotalOrdersOfDelivery(ctx context.Context, deliveryId string, keyword string, query *pagable.Query) (int, error)
	TotalSearchOrderByStoreID(ctx context.Context, storeId string, keyword string) (int, error)
	//custom_entity - admin
	GetTotalOrderInSystemInDay(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInHours, error)
	GetTotalOrderInSystemInMonth(ctx context.Context, month int, year int) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYear(ctx context.Context, year int, count int) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYear(ctx context.Context, date string) ([]custom_entity.OrderCommissionDetail, error)
	ListOfProductSelledOnMonth(ctx context.Context, date string, count int) ([]custom_entity.TopOfProductSold, error)

	//custom_entity - store
	GetTotalOrderInSystemInMonthOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYearOfStore(ctx context.Context, year int, storeId string, count int) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYearOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.OrderCommissionDetail, error)
	ListOfProductSelledOnMonthStore(ctx context.Context, date string, count int, storeId string) ([]custom_entity.TopOfProductSold, error)

	UserCountingOrder(ctx context.Context, userId string) (int, error)
	StoreCountingOrder(ctx context.Context, storeId string) (int, error)
	DeliveryCountingOrder(ctx context.Context, deliveryId string) (int, error)
	AdminCountingOrder(ctx context.Context) (int, error)
}
