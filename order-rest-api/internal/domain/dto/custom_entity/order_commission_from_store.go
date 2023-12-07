package custom_entity

type OrderCommissionDetail struct {
	Month         int `json:"month"`
	TotalReceived int `json:"total_received"`
	TotalFee      int `json:"total_fee"`
	TotalOrders   int `json:"total_orders"`
}

type AmountItemOfStoreInOrder struct {
	StoreId string `json:"store_id"`
	Total   int    `json:"total"`
}
