package pagable

type PageableQuery struct {
	Page string `json:"page"`
	Size string `json:"size"`
}

type Query struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ListResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
