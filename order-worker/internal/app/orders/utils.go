package orders

import (
	"order-worker/internal/domain/dto/order"
	"order-worker/internal/infrastructure/adapter/productserv/dto"
)

func MappingOrderItemForReduce(request []order.OrderItemsCache) []dto.ReduceItem {
	var items []dto.ReduceItem
	for _, i := range request {
		product := dto.ReduceItem{
			ProductId: i.ProductItem.ProductID,
			OptionId:  i.ProductItem.OptionID,
			Quantity:  i.ProductItem.Quantity,
		}
		items = append(items, product)
	}
	return items
}
