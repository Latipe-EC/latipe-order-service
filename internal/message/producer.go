package message

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/domain/entities/order"
)

type ConsumerOrderMessage struct {
	config     *config.Config
	connection *amqp.Connection
	orderRepo  order.Repository
}

func NewConsumerOrderMessage(config *config.Config, connect *amqp.Connection, repository order.Repository) *ConsumerOrderMessage {
	return &ConsumerOrderMessage{
		config:     config,
		connection: connect,
		orderRepo:  repository,
	}
}

func (mq ConsumerOrderMessage) ListenMessageQueue() {
	channel, err := mq.connection.Channel()
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
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()
	fmt.Println("Waiting for messages...")
	<-forever
}
