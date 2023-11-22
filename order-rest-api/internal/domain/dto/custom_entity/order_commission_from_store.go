package custom_entity

type OrderCommissionDetail struct {
	Month         int `json:"month"`
	Amount        int `json:"amount"`
	TotalReceived int `json:"total_received"`
	TotalFee      int `json:"total_fee"`
	Count         int `json:"count"`
}

type AmountItemOfStoreInOrder struct {
	StoreId string `json:"store_id"`
	Total   int    `json:"total"`
}
