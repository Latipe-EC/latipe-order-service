package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type OrderCommissionDetailRequest struct {
	Date    string `json:"date"`
	Count   int    `json:"count"`
	StoreId string
}

type OrderCommissionDetailResponse struct {
	StoreID string                                `json:"store_id,omitempty"`
	Items   []custom_entity.OrderCommissionDetail `json:"items"`
}
