package order

type CheckUserOrderRequest struct {
	Header    BaseHeader
	ProductId string `query:"product_id" validate:"required"`
	UserId    string `validate:"required"`
}

type CheckUserOrderResponse struct {
	IsPurchased bool     `json:"is_purchased"`
	Orders      []string `json:"orders,omitempty"`
}
