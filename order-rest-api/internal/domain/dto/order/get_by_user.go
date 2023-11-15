package order

import "order-rest-api/pkg/util/pagable"

type GetByUserIdRequest struct {
	BaseHeader BaseHeader
	UserId     string
	Query      *pagable.Query
}

type GetByUserIdResponse struct {
	pagable.ListResponse
}
