package worker

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"order-worker/config"
	"order-worker/internal/app/orders"
	order2 "order-worker/internal/domain/dto/order"
	"order-worker/internal/domain/entities/order"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerOrderMessage struct {
	config       *config.Config
	orderUsecase orders.Usecase
}

func NewConsumerOrderMessage(config *config.Config, orderService orders.Usecase) *ConsumerOrderMessage {
	return &ConsumerOrderMessage{
		config:       config,
		orderUsecase: orderService,
	}
}

func (mq ConsumerOrderMessage) ListenOrderEventQueue(wg *sync.WaitGroup) {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Info("Comsumer has been connected")

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// Khai báo một Exchange loại "direct"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.OrderEvent.Exchange, // Tên Exchange
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare exchange: %v", err)
	}

	// Tạo hàng đợi
	_, err = channel.QueueDeclare(
		mq.config.RabbitMQ.OrderEvent.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("cannot declare queue: %v", err)
	}

	err = channel.QueueBind(
		mq.config.RabbitMQ.OrderEvent.Queue,
		mq.config.RabbitMQ.OrderEvent.RoutingKey,
		mq.config.RabbitMQ.OrderEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		mq.config.RabbitMQ.OrderEvent.Queue, // queue
		mq.config.RabbitMQ.ConsumerName,     // consumer
		true,                                // auto ack
		false,                               // exclusive
		false,                               // no local
		false,                               // no wait
		nil,                                 //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof("received order message from: %s", msg.RoutingKey)
		if err := mq.orderHandler(msg); err != nil {
			log.Infof("The order creation failed cause %s", err)
		}

	}

	log.Infof("message queue has started")
	log.Infof("waiting for messages...")

}

func (mq ConsumerOrderMessage) orderHandler(msg amqp.Delivery) error {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := order2.OrderMessage{}

	if err := json.Unmarshal(msg.Body, &message); err != nil {
		log.Infof("Parse message to order failed cause: %s", err)
		return err
	}

	switch message.Status {

	case order.ORDER_SYSTEM_PROCESS:
		err := mq.orderUsecase.CreateOrderTransaction(ctx, &message)
		if err != nil {
			return err
		}
	case order.ORDER_CANCEL:
		err := mq.orderUsecase.RollBackCancelOrder(ctx, message.OrderUUID)
		if err != nil {
			return err
		}
	case order.ORDER_REFUND:
		err := mq.orderUsecase.RollBackCancelOrder(ctx, message.OrderUUID)
		if err != nil {
			return err
		}
	case order.ORDER_SHIPPING_FINISH:
		return nil
	}

	endTime := time.Now()
	log.Infof("The order is processed successfully - duration:%v", endTime.Sub(startTime))
	return nil
}
