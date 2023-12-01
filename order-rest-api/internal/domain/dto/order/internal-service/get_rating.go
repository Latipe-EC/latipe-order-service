package internal_service

type GetOrderRatingItemRequest struct {
	OrderUUID string `params:"id" json:"order_uuid"`
}

type GetOrderRatingItemResponse struct {
	OrderUUID  string             `json:"order_uuid"`
	OrderItems []OrderRatingItems `json:"order_items"`
}

type OrderRatingItems struct {
	ItemId    string `json:"item_id"`
	ProductId string `json:"product_id" `
	OptionId  string `json:"option_id"`
	RatingID  string `json:"rating_id,omitempty"`
}
