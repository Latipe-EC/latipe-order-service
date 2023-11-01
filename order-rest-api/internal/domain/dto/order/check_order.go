package order

type CheckUserOrderRequest struct {
	Header    BaseHeader
	ProductId string `json:"product_id"`
	UserId    string
}

type CheckUserOrderResponse struct {
	IsPurchased bool     `json:"is_purchased"`
	Orders      []string `json:"orders,omitempty"`
}
