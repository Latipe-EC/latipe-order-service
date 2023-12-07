package statistic

import "order-rest-api/internal/domain/dto/custom_entity"

type AdminTotalOrderInDayRequest struct {
	Date string `query:"date"`
}
type AdminTotalOrderInDayResponse struct {
	FilterDate string                                    `json:"filter_date,omitempty"`
	Items      []custom_entity.TotalOrderInSystemInHours `json:"items"`
}

type AdminTotalOrderInMonthRequest struct {
	Date  string `query:"date"`
	Count int    `json:"count"`
}

type AdminTotalOrderInMonthResponse struct {
	FilterDate string                                  `json:"filter_date,omitempty"`
	Items      []custom_entity.TotalOrderInSystemInDay `json:"items"`
}

type AdminGetTotalOrderInYearRequest struct {
	Year int `query:"year"`
}

type AdminGetTotalOrderInYearResponse struct {
	FilterDate string                                    `json:"filter_date,omitempty"`
	Items      []custom_entity.TotalOrderInSystemInMonth `json:"items"`
}
