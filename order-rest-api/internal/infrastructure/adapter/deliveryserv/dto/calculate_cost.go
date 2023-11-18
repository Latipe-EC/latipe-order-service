package dto

const shippingCostUrl = "/api/v1/delivery/shipping/order"

type GetShippingCostRequest struct {
	SrcCode    []string `json:"src_code"`
	DestCode   string   `json:"dest_code"`
	DeliveryId string   `json:"delivery_id"`
}

type GetShippingCostResponse struct {
	ReceiveDate  string `json:"receive_date"`
	DeliveryId   string `json:"delivery_id"`
	DeliveryName string `json:"delivery_name"`
	Cost         int    `json:"cost"`
}

func (GetShippingCostRequest) URL() string {
	return shippingCostUrl
}
