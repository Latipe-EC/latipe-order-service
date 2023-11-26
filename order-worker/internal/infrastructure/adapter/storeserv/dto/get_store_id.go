package dto

const getStoreId = "/api/v1/stores/"

type GetStoreByIdRequest struct {
	BaseHeader
	StoreID string `json:"store_id"`
}

type GetStoreByIdResponse struct {
	Id          string  `json:"id"`
	IsDeleted   bool    `json:"isDeleted"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Logo        string  `json:"logo"`
	OwnerId     string  `json:"ownerId"`
	Cover       string  `json:"cover"`
	FeePerOrder float64 `json:"feePerOrder"`
}

func (GetStoreByIdRequest) URL() string {
	return getStoreId
}
