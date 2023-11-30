package dto

const rollbackURL = "/api/v1/vouchers/rollback"

type RollbackVoucherRequest struct {
	Vouchers []string `json:"vouchers"`
	AuthorizationHeader
}

func (RollbackVoucherRequest) URL() string {
	return url
}
