package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type OrderCommissionDetailRequest struct {
	Date    string `json:"date" query:"date"`
	Count   int    `json:"count" query:"count"`
	StoreId string
}

type OrderCommissionDetailResponse struct {
	StoreID    string                                `json:"store_id,omitempty"`
	FilterDate string                                `json:"filter_date,omitempty"`
	Items      []custom_entity.OrderCommissionDetail `json:"items"`
}
