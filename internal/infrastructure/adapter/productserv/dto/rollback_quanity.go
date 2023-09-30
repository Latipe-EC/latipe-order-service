package dto

const rollbackQuantityUrl = "/api/v1/products/quantity"

type RollbackQuantityRequest struct {
	Items []RollBackItem `json:"items"`
}

type RollBackItem struct {
	ProductId string `json:"productId"`
	OptionId  int    `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type RollbackQuantityResponse struct {
	Message string `json:"message"`
}

func (RollbackQuantityRequest) URL() string {
	return rollbackQuantityUrl
}
