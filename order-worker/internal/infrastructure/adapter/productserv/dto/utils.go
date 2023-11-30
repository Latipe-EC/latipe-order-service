package dto

import "order-worker/internal/domain/dto/order"

func MappingReduceProduct(item []order.OrderItemsCache) []ReduceItem {
	var data []ReduceItem
	for _, i := range item {

		newItem := ReduceItem{
			ProductId: i.ProductItem.ProductID,
			OptionId:  i.ProductItem.OptionID,
			Quantity:  i.ProductItem.Quantity,
		}
		data = append(data, newItem)
	}

	return data
}

func MappingeRollbackProduct(item []order.OrderItemsCache) []ReduceItem {
	var data []ReduceItem
	for _, i := range item {

		newItem := ReduceItem{
			ProductId: i.ProductItem.ProductID,
			OptionId:  i.ProductItem.OptionID,
			Quantity:  -i.ProductItem.Quantity,
		}
		data = append(data, newItem)
	}

	return data
}
