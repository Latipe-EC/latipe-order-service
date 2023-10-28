package order

type CreateOrderRequest struct {
	Header      BaseHeader
	UserRequest UserRequest

	PaymentMethod int          `json:"payment_method" validate:"required"`
	VoucherCode   []Voucher    `json:"vouchers"`
	Address       OrderAddress `json:"address" validate:"required"`
	Delivery      Delivery     `json:"delivery" validate:"required"`
	OrderItems    []OrderItems `json:"order_items" validate:"required"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type Voucher struct {
	Code string `json:"code"`
	Type int    `json:"type"`
}

type OrderItems struct {
	CartItemId string `json:"cart_item_id"`
	ProductId  string `json:"product_id" validate:"required"`
	OptionId   string `json:"option_id"`
	Quantity   int    `json:"quantity" validate:"required"`
	Price      int    `json:"price" validate:"required"`
}

type OrderAddress struct {
	AddressId string `json:"address_id" validate:"required"`
}
type Delivery struct {
	DeliveryId string `json:"delivery_id" validate:"required"`
}
