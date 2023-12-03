package internal_service

type GetOrderRatingItemRequest struct {
	ItemID string `params:"id" json:"item_id"`
}

type GetOrderRatingItemResponse struct {
	RatingId string `json:"rating_id"`
}
