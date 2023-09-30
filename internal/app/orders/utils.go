package orders

import (
	"order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/infrastructure/adapter/productserv/dto"
)

func MappingOrderItemForReduce(request *order.CreateOrderRequest) []dto.ReduceItem {
	var items []dto.ReduceItem
	for _, i := range request.OrderItems {
		product := dto.ReduceItem{
			ProductId: i.ProductId,
			OptionId:  i.OptionId,
			Quantity:  i.Quantity,
		}
		items = append(items, product)
	}
	return items
}
