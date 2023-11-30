package dto

const reduceProductUrl = "/api/v1/products/update-quantity"

type ReduceProductRequest struct {
	Items []ReduceItem `json:"items"`
}

type ReduceItem struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type ReduceProductResponse struct {
	Message string `json:"message"`
}

func (ReduceProductRequest) URL() string {
	return reduceProductUrl
}
