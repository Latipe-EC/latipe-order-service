package order

type CreateOrderRequest struct {
	Header        BaseHeader
	UserRequest   UserRequest  `json:"user_request"`
	Amount        int          `json:"amount" validate:"required"`
	Discount      int          `json:"discount" validate:"required"`
	Total         int          `json:"total" validate:"required"`
	PaymentMethod int          `json:"payment_method" validate:"required"`
	VoucherCode   string       `json:"voucher_code" validate:"required"`
	CreateAt      string       `json:"create_at"`
	Address       OrderAddress `json:"address" validate:"required"`
	OrderItems    []OrderItems `json:"order_items" validate:"required"`
}
type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OrderItems struct {
	ProductId string `json:"product_id" validate:"required"`
	OptionId  int    `json:"option_id"`
	Quantity  int    `json:"quantity" validate:"required"`
	Price     int    `json:"price" validate:"required"`
}
type OrderAddress struct {
	AddressId       string `json:"address_id" validate:"required"`
	ShippingName    string `json:"shipping_name" validate:"required"`
	ShippingPhone   string `json:"shipping_phone" validate:"required"`
	ShippingAddress string `json:"shipping_address" validate:"required"`
}
type Delivery struct {
	DeliveryId    string `json:"delivery_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Cost          int    `json:"cost" validate:"required"`
	ReceivingDate string `json:"receiving_date" validate:"required"`
}
