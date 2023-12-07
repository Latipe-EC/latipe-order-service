package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type ListOfProductSoldRequest struct {
	Date    string `query:"date"`
	Count   int    `query:"count"`
	StoreId string
}

type ListOfProductSoldResponse struct {
	StoreID    string                           `json:"store_id,omitempty"`
	FilterDate string                           `json:"filter_date,omitempty"`
	Items      []custom_entity.TopOfProductSold `json:"items"`
}
