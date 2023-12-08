package dto

import (
	"order-worker/internal/domain/dto/order"
	entity "order-worker/internal/domain/entities/order"
)

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

func MappingDAORollbackProduct(item []*entity.OrderItem) []ReduceItem {
	var data []ReduceItem
	for _, i := range item {

		newItem := ReduceItem{
			ProductId: i.ProductID,
			OptionId:  i.OptionID,
			Quantity:  -i.Quantity,
		}
		data = append(data, newItem)
	}

	return data
}
