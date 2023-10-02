package dto

const orderProductUrl = "/api/v1/products/check-in-stock"

type OrderProductRequest struct {
	Items []ValidateItems
}

type OrderProductResponse struct {
	Products   []Product `json:"products"`
	TotalPrice int       `json:"totalPrice"`
}

type ValidateItems struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type Product struct {
	ProductId        string `json:"id"`
	Name             string `json:"name"`
	Image            string `json:"image"`
	Code             string `json:"code"`
	Sku              string `json:"sku"`
	Quantity         int    `json:"quantity"`
	Price            int    `json:"price"`
	NameOption       string `json:"nameOption"`
	PromotionalPrice int    `json:"promotionalPrice"`
}

func (OrderProductRequest) URL() string {
	return orderProductUrl
}
