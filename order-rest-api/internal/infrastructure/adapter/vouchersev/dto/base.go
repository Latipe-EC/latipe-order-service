package dto

type AuthorizationHeader struct {
	BearerToken string `reqHeader:"Authorization" json:"bearer_token"`
}

const (
	FREE_SHIP      = 1
	DISCOUNT_ORDER = 2

	PENDING   = 0
	ACTIVE    = 1
	IN_ACTIVE = 2
)
