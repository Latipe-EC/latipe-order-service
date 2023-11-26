package dto

const getStoreId = "/api/v1/stores/"

type GetStoreByIdRequest struct {
	BaseHeader
	StoreID string `json:"store_id"`
}

type GetStoreByIdResponse struct {
	StoreID string  `json:"store_id"`
	Fee     float64 `json:"fee"`
}

func (GetStoreByIdRequest) URL() string {
	return getStoreId
}
