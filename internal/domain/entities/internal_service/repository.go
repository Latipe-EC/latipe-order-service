package internal_service

type Repository interface {
	FindByClient(clientId string, secretKey string) (*InternalService, error)
	Save(service InternalService) error
	Update(service InternalService) error
}
