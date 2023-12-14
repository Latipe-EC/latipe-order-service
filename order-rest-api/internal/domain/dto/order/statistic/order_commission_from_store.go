package statistic

type OrderCommissionDetailRequest struct {
	Date    string `json:"date" query:"date"`
	StoreId string
}

type OrderCommissionDetailResponse struct {
	StoreID    string      `json:"store_id,omitempty"`
	FilterDate string      `json:"filter_date,omitempty"`
	Items      interface{} `json:"items"`
}
