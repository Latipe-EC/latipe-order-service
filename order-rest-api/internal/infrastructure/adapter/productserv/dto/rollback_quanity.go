package dto

const rollbackQuantityUrl = "/api/v1/products/update-quantity"

type RollbackQuantityRequest struct {
	Items []RollBackItem
}

type RollBackItem struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type RollbackQuantityResponse struct {
	Message string `json:"message"`
}

func (RollbackQuantityRequest) URL() string {
	return rollbackQuantityUrl
}
