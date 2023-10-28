package message

import (
	"context"
	"encoding/json"
	"log"
	"order-worker/config"
	"order-worker/internal/app/orders"
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

func (mq ConsumerOrderMessage) ListenMessageQueue() {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] Comsumer has been connected", "INFO")

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		mq.config.RabbitMQ.OrderQueue,   // queue
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

	var orderCacheKey string
	if err := json.Unmarshal(msg.Body, &orderCacheKey); err != nil {
		log.Printf("[%s] Parse message to order failed cause: %s", "ERROR", err)
		return err
	}

	err := mq.orderUsecase.CreateOrder(ctx, orderCacheKey)
	if err != nil {
		return err
	}

	log.Printf("[%s] The order created successfully: %s", "INFO", msg.RoutingKey)
	return nil
}
