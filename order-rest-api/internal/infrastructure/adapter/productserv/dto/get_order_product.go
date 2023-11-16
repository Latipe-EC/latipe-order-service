package dto

const orderProductUrl = "/api/v1/products/check-in-stock"

type OrderProductRequest struct {
	Items []ValidateItems
}

type ValidateItems struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type OrderProductResponse struct {
	Products           []Product `json:"products"`
	TotalPrice         int       `json:"totalPrice"`
	StoreProvinceCodes []string  `json:"storeProvinceCodes"`
}

type Product struct {
	ProductId        string  `json:"productId"`
	Name             string  `json:"name"`
	Quantity         int     `json:"quantity"`
	Image            string  `json:"image"`
	Price            float64 `json:"price"`
	PromotionalPrice float64 `json:"promotionalPrice"`
	OptionId         string  `json:"optionId"`
	NameOption       string  `json:"nameOption"`
	StoreId          string  `json:"storeId"`
	TotalPrice       float64 `json:"totalPrice"`
}

func (OrderProductRequest) URL() string {
	return orderProductUrl
}
