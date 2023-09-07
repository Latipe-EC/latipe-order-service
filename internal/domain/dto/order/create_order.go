package order

type CreateOrderRequest struct {
	Header        BaseHeader
	Amount        int          `json:"amount"`
	Discount      int          `json:"discount"`
	Total         int          `json:"total"`
	PaymentMethod int          `json:"payment_method"`
	VoucherCode   string       `json:"voucher_code"`
	CreateAt      string       `json:"create_at"`
	Address       OrderAddress `json:"address"`
	OrderItems    []OrderItems `json:"order_items"`
}

type OrderItems struct {
	ProductId string `json:"product_id"`
	OptionId  string `json:"option_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}
type OrderAddress struct {
	AddressId     string `json:"address_id"`
	AddressDetail string `json:"address_detail"`
}
type Delivery struct {
	DeliveryId    string `json:"delivery_id"`
	Name          string `json:"name"`
	Cost          int    `json:"cost"`
	ReceivingDate string `json:"receiving_date"`
}
