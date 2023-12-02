package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type AdminTotalOrderInDayRequest struct {
	Date string `query:"date"`
}
type AdminTotalOrderInDayResponse struct {
	Items []custom_entity.TotalOrderInSystemInHours `json:"items"`
}

type AdminTotalOrderInMonthRequest struct {
	Month  int `json:"month"`
	Amount int `json:"amount"`
	Count  int `json:"count"`
}

type AdminTotalOrderInMonthResponse struct {
	Items []custom_entity.TotalOrderInSystemInDay `json:"items"`
}

type AdminGetTotalOrderInYearRequest struct {
	Year int `query:"year"`
}

type AdminGetTotalOrderInYearResponse struct {
	Items []custom_entity.TotalOrderInSystemInMonth `json:"items"`
}
