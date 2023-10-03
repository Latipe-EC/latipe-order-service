package dto

const orderProductUrl = "/api/v1/products/check-in-stock"

type OrderProductRequest struct {
	Items []ProductItem `json:"products"`
}

type OrderProductResponse struct {
	Products   []Products `json:"products"`
	TotalPrice int        `json:"totalPrice"`
}

type ProductItem struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type Products struct {
	ProductId        string `json:"productId"`
	Name             string `json:"name"`
	Quantity         int    `json:"quantity"`
	Price            int    `json:"price"`
	Discount         int    `json:"discount"`
	OptionId         string `json:"optionId"`
	NameOption       string `json:"nameOption"`
	StoreId          string `json:"storeId"`
	PromotionalPrice int    `json:"promotionalPrice"`
}

func (OrderProductRequest) URL() string {
	return orderProductUrl
}
