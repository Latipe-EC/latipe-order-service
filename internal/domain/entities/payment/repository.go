package payment

import "log"

type PaymentRepository interface {
	FindById(id int) (*PaymentLog, error)
	FindByOrderId(id int) (*PaymentLog, error)
	Save(log PaymentLog) error
	Update(logger log.Logger) error
}
