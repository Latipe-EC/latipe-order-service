package message

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/domain/entities/order"
)

type ConsumerOrderMessage struct {
	config    *config.Config
	orderRepo order.Repository
}

func NewConsumerOrderMessage(config *config.Config, repository order.Repository) *ConsumerOrderMessage {
	return &ConsumerOrderMessage{
		config:    config,
		orderRepo: repository,
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
	}

	log.Printf("[%s] message queue has started", "INFO")
	log.Printf("[%s] waiting for messages...", "INFO")

}
