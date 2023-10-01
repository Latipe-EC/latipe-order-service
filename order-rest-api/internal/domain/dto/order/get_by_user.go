package order

import "order-rest-api/pkg/util/pagable"

type GetByUserIdRequest struct {
	BaseHeader BaseHeader
	UserId     int
	Query      *pagable.Query
}

type GetByUserIdResponse struct {
	pagable.ListResponse
}
