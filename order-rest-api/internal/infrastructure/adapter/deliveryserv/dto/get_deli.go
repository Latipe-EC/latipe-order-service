package dto

import "time"

const getDeliURL = "/api/v1/delivery/validate"

type GetDeliveryByTokenRequest struct {
	BearerToken string `json:"bearer_token"`
}

type GetDeliveryByTokenResponse struct {
	ID           string    `json:"_id" `
	DeliveryName string    `json:"delivery_name,omitempty"`
	DeliveryCode string    `json:"delivery_code,omitempty"`
	BaseCost     int       `json:"base_cost,omitempty" `
	Description  string    `json:"description,omitempty" `
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	IsActive     bool      `json:"is_active"`
}

func (GetDeliveryByTokenRequest) URL() string {
	return getDeliURL
}
