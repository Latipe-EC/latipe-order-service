package order

type BaseHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}

type OrderResponse struct {
	Amount        int          `json:"amount"`
	Discount      int          `json:"discount"`
	Total         int          `json:"total"`
	Status        int          `json:"status"`
	PaymentMethod int          `json:"payment_method"`
	CreateAt      string       `json:"create_at"`
	UpdateAt      string       `json:"update_at"`
	Address       OrderAddress `json:"address"`
	OrderItems    []OrderItems `json:"order_items"`
}
