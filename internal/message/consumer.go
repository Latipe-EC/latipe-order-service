package message

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/app/orders"
	dto "order-service-rest-api/internal/domain/dto/order"
	"order-service-rest-api/internal/domain/entities/order"
	payment "order-service-rest-api/internal/domain/entities/payment"
	"order-service-rest-api/pkg/util/mapper"
	"time"
)

type ConsumerOrderMessage struct {
	config       *config.Config
	orderRepo    order.Repository
	orderUsecase orders.Usecase
}

func NewConsumerOrderMessage(config *config.Config, orderRepos order.Repository, orderService orders.Usecase) *ConsumerOrderMessage {
	return &ConsumerOrderMessage{
		config:       config,
		orderRepo:    orderRepos,
		orderUsecase: orderService,
	}
}

func (mq ConsumerOrderMessage) ListenMessageQueue() {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] Comsumer has been connected", "INFO")

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		mq.config.RabbitMQ.Queue,        // queue
		mq.config.RabbitMQ.ConsumerName, // consumer
		true,                            // auto ack
		false,                           // exclusive
		false,                           // no local
		false,                           // no wait
		nil,                             //args
	)
	if err != nil {
		panic(err)
	}

	// handle consumed messages from queue
	for msg := range msgs {
		log.Printf("[%s] received order message from: %s", "INFO", msg.RoutingKey)

		if err := mq.orderHandler(msg); err != nil {
			log.Printf("[%s] The order creation failed cause %s", "ERROR", err)
		}

	}

	log.Printf("[%s] message queue has started", "INFO")
	log.Printf("[%s] waiting for messages...", "INFO")

}

func (mq ConsumerOrderMessage) orderHandler(msg amqp.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orderDTO dto.CreateOrderRequest
	if err := json.Unmarshal(msg.Body, &orderDTO); err != nil {
		log.Printf("[%s] Parse message to order failed cause: %s", "ERROR", err)
		return err
	}

	orderDAO := order.Order{}

	if err := mapper.BindingStruct(orderDTO, &orderDAO); err != nil {
		log.Printf("[%s] Mapping value from dto to dao failed cause: %s", "ERROR", err)
		return err
	}

	//create items
	var orderItems []*order.OrderItem
	for _, item := range orderDTO.OrderItems {
		i := order.OrderItem{
			ProductID: item.ProductId,
			SellerID:  0,
			OptionID:  item.OptionId,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Order:     &orderDAO,
		}
		orderItems = append(orderItems, &i)
	}
	orderDAO.OrderItem = orderItems

	//create log
	var logs []*order.OrderStatusLog
	orderLog := order.OrderStatusLog{
		Order:        &orderDAO,
		Message:      "order created",
		StatusChange: order.ORDER_CREATED,
	}
	orderDAO.OrderStatusLog = append(logs, &orderLog)

	//create payment
	paymentLog := payment.PaymentLog{
		PaymentTransaction: uuid.New().String(),
		PaymentType:        orderDTO.PaymentMethod,
		Total:              0,
		ThirdPartyLog:      "",
	}
	orderDAO.PaymentLog = &paymentLog

	err := mq.orderRepo.Save(orderDAO)
	if err != nil {
		if err := mq.orderUsecase.RollBackQuantity(ctx, &orderDTO); err != nil {
			log.Printf("[%s] Rollback product quantity was failed : %v", "ERROR", err)
			return err
		}

		log.Printf("[%s] The order created failed : %s", "ERROR")
		return err
	}

	log.Printf("[%s] The order created successfully: %s", "INFO", msg.RoutingKey)
	return nil
}
