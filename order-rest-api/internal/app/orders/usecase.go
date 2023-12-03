package orders

import (
	"context"
	orderDTO "order-rest-api/internal/domain/dto/order"
	"order-rest-api/internal/domain/dto/order/delivery"
	internalDTO "order-rest-api/internal/domain/dto/order/internal-service"
	"order-rest-api/internal/domain/dto/order/statistic"
	"order-rest-api/internal/domain/dto/order/store"
)

type Usecase interface {
	//admin
	GetOrderById(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error)
	UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error
	GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error)
	CheckProductPurchased(ctx context.Context, dto *orderDTO.CheckUserOrderRequest) (*orderDTO.CheckUserOrderResponse, error)
	UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error
	AdminCountingOrderAmount(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	// internal service
	InternalGetRatingID(ctx context.Context, dto *internalDTO.GetOrderRatingItemRequest) (*internalDTO.GetOrderRatingItemResponse, error)

	//user
	ProcessCacheOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error)
	GetOrderByUUID(ctx context.Context, dto *orderDTO.GetOrderByUUIDRequest) (*orderDTO.GetOrderResponse, error)
	CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error)
	UserCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	//store
	GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error)
	ViewDetailStoreOrder(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error)
	UpdateOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error)
	StoreCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	//deli
	DeliveryUpdateStatusOrder(ctx context.Context, dto delivery.UpdateOrderStatusRequest) (*delivery.UpdateOrderStatusResponse, error)
	GetOrdersOfDelivery(ctx context.Context, dto *delivery.GetOrderListRequest) (*delivery.GetOrderListResponse, error)
	DeliveryCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	//custom_entity - admin
	AdminGetTotalOrderInSystemInDay(ctx context.Context, dto *statistic.AdminTotalOrderInDayRequest) (*statistic.AdminTotalOrderInDayResponse, error)
	AdminGetTotalOrderInSystemInMonth(ctx context.Context, dto *statistic.AdminTotalOrderInMonthRequest) (*statistic.AdminTotalOrderInMonthResponse, error)
	AdminGetTotalOrderInSystemInYear(ctx context.Context, dto *statistic.AdminGetTotalOrderInYearRequest) (*statistic.AdminGetTotalOrderInYearResponse, error)
	AdminGetTotalCommissionOrderInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error)
	AdminListOfProductSoldOnMonth(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error)

	//custom_entity - store
	GetTotalOrderInMonthOfStore(ctx context.Context, dto *statistic.GetTotalStoreOrderInMonthRequest) (statistic.GetTotalOrderInMonthResponse, error)
	GetTotalOrderInYearOfStore(ctx context.Context, dto *statistic.GetTotalOrderInYearOfStoreRequest) (*statistic.GetTotalOrderInYearOfStoreResponse, error)
	GetTotalStoreCommissionInYear(ctx context.Context, dto *statistic.OrderCommissionDetailRequest) (*statistic.OrderCommissionDetailResponse, error)
	ListOfProductSoldOnMonthStore(ctx context.Context, dto *statistic.ListOfProductSoldRequest) (*statistic.ListOfProductSoldResponse, error)
}
