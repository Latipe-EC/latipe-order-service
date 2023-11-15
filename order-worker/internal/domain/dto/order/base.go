package order

type BaseHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}

type OrderResponse struct {
	OrderKey      string            `json:"order_key"`
	Amount        int               `json:"amount"`
	Discount      int               `json:"discount"`
	Total         int               `json:"total"`
	Status        int               `json:"status"`
	PaymentMethod int               `json:"payment_method"`
	VoucherCode   string            `json:"voucher_code"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	Address       OrderAddress      `json:"address"`
	OrderItems    []OrderItemsCache `json:"order_items"`
}
