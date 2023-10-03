package orders

import (
	"order-worker/internal/domain/dto/order"
	"order-worker/internal/infrastructure/adapter/productserv/dto"
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

func MappingOrderItemForRollback(request *order.CreateOrderRequest) []dto.RollbackItem {
	var items []dto.RollbackItem
	for _, i := range request.OrderItems {
		product := dto.RollbackItem{
			ProductId: i.ProductId,
			OptionId:  i.OptionId,
			Quantity:  i.Quantity,
		}
		items = append(items, product)
	}
	return items
}

func MappingOrderItemToValidateItems(items []order.OrderItems) []dto.ProductItem {
	var results []dto.ProductItem
	for _, i := range items {
		reqItem := dto.ProductItem{
			ProductId: i.ProductId,
			OptionId:  i.OptionId,
			Quantity:  i.Quantity,
		}
		results = append(results, reqItem)
	}
	return results
}
