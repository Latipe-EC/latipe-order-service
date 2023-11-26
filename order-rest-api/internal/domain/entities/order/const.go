package order

const (
	ORDER_SYSTEM_PROCESS = 0
	ORDER_CREATED        = 1
	ORDER_PENDING        = 2
	ORDER_DELIVERY       = 3

	ORDER_SHIPPING_FINISH = 4
	ORDER_COMPLETED       = 5
	ORDER_REFUND          = 6
	ORDER_CANCEL          = -1
)

const (
	PAID_COD     = 1
	PAID_VIA_3RD = 2
)
