package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type GetTotalStoreOrderInMonthRequest struct {
	Date    string `query:"date"`
	StoreId string
}

type GetTotalOrderInMonthResponse struct {
	Items []custom_entity.TotalOrderInSystemInDay `json:"items"`
}

type GetTotalOrderInYearOfStoreRequest struct {
	Year    int `query:"year"`
	StoreID string
}

type GetTotalOrderInYearOfStoreResponse struct {
	Items []custom_entity.TotalOrderInSystemInMonth `json:"items"`
}
