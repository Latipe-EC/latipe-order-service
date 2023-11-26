package dto

type StoreBillingMessage struct {
	StoreID        string `json:"storeId"`
	OrderUUID      string `json:"orderUuid"`
	AmountReceived int    `json:"amountReceived"`
	SystemFee      int    `json:"systemFee"`
}
