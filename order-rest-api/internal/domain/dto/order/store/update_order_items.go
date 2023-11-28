package store

type UpdateOrderItemRequest struct {
	OrderUUID string `params:"id" validate:"required"`
	ItemID    int    `json:"item_id"`
	StoreId   string
}

type UpdateOrderItemResponse struct {
	OrderUUID string `json:"order_uuid"`
	ItemID    int    `json:"item_id"`
	Status    int    `json:"status"`
}
