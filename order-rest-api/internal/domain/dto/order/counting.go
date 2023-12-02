package order

type CountingOrderAmountRequest struct {
	Role    string
	OwnerID string
}

type CountingOrderAmountResponse struct {
	Count int `json:"count"`
}
