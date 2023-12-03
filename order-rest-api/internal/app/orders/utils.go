package orders

import (
	"order-rest-api/internal/domain/dto/order"
	enitites "order-rest-api/internal/domain/entities/order"
	"order-rest-api/internal/infrastructure/adapter/productserv/dto"
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

func CheckStoreHaveOrder(entities enitites.Order, storeId string) bool {
	for _, i := range entities.OrderItem {
		if i.StoreID == storeId {
			return true
		}
	}
	return false
}

func MappingOrderItemToGetInfo(request *order.CreateOrderRequest) []dto.ValidateItems {
	var items []dto.ValidateItems
	for _, i := range request.OrderItems {
		product := dto.ValidateItems{
			ProductId: i.ProductId,
			OptionId:  i.OptionId,
			Quantity:  i.Quantity,
		}
		items = append(items, product)
	}
	return items
}

func deleteItems(slice []order.OrderItems, index int) []order.OrderItems {
	return append(slice[:index], slice[index+1:]...)
}
