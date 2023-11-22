package statistic

type TopOfProductSold struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Total       int    `json:"total"`
}
