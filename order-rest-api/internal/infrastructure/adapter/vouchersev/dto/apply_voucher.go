package dto

const url = "/api/v1/vouchers/checking"

type CheckingVoucherRequest struct {
	OrderTotalAmount int      `json:"order_total_amount"`
	PaymentMethod    int      `json:"payment_method"`
	UserId           string   `json:"user_id"`
	Vouchers         []string `json:"vouchers"`
	AuthorizationHeader
}

type UseVoucherResponse struct {
	IsSuccess bool                `json:"is_success"`
	Items     []VoucherRespDetail `json:"items"`
}

type VoucherRespDetail struct {
	ID              string         `json:"_id,omitempty"`
	VoucherCode     string         `json:"voucher_code"`
	VoucherType     int            `json:"voucher_type"`
	OwnerVoucher    string         `json:"owner_voucher,omitempty"`
	DiscountPercent float64        `json:"discount_percent,omitempty"`
	DiscountValue   int            `json:"discount_value,omitempty"`
	VoucherRequire  VoucherReqResp `json:"voucher_require,omitempty"`
	Status          int            `json:"status"`
}

type VoucherReqResp struct {
	MinRequire          int64  `json:"min_require"`
	MemberType          int    `json:"member_type,omitempty"`
	PaymentMethod       int    `json:"payment_method,omitempty"`
	RequiredOwnerProdId string `json:"required_owner_prod_id,omitempty"`
}

func (CheckingVoucherRequest) URL() string {
	return url
}
