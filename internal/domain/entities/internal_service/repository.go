package internal_service

type Repository interface {
	FindByKeyAndValue(key string, value string) (*InternalService, error)
	Save(service InternalService) error
	Update(service InternalService) error
}
