package order

const (
	ORDER_SYSTEM_PROCESS = 0
	ORDER_CREATED        = 1
	ORDER_PENDING        = 2
	ORDER_DELIVERY       = 3

	ORDER_SHIPPING_FINISH = 4
	ORDER_COMPLETED       = 5
	ORDER_REFUND          = 6
	ORDER_CANCEL          = 7
	ORDER_FAILED          = -1
)

const (
	PAYMENT_COD        = 1
	PAYMENT_VIA_PAYPAL = 2
	PAYMENT_VIA_WALLET = 3
)

const (
	CMS_DONE    = 1
	CMS_PENDING = 0
	CMS_CANCEL  = -1
)
