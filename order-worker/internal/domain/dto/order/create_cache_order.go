package order

type OrderCacheData struct {
	Header        BaseHeader
	UserRequest   UserRequest       `json:"user_request"`
	Amount        int               `json:"amount" validate:"required"`
	Discount      int               `json:"discount" validate:"required"`
	SubTotal      int               `json:"sub_total" validate:"required"`
	PaymentMethod int               `json:"payment_method" validate:"required"`
	Vouchers      []Voucher         `json:"vouchers"`
	Address       OrderAddress      `json:"address" validate:"required"`
	Delivery      Delivery          `json:"delivery" validate:"required"`
	OrderItems    []OrderItemsCache `json:"order_items" validate:"required"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type Voucher struct {
	Code string `json:"code"`
	Type int    `json:"type"`
}

type OrderItemsCache struct {
	CartItemId  string      `json:"cart_item_id"`
	ProductItem ProductItem `json:"product_item"`
}

type ProductItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	StoreID     string `json:"store_id"`
	OptionID    string `json:"option_id" `
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	NetPrice    int    `json:"net_price"`
}

type OrderAddress struct {
	AddressId       string `json:"address_id"`
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
