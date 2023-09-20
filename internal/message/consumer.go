package message

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"order-service-rest-api/config"
	"time"
)

type ProducerOrderMessage struct {
	config     *config.Config
	connection *amqp.Connection
}

func NewProducerOrderMessage(config *config.Config, connect *amqp.Connection) *ProducerOrderMessage {
	return &ProducerOrderMessage{
		config:     config,
		connection: connect,
	}
}

func (mq *ConsumerOrderMessage) SendMessage(request interface{}) error {
	ch, err := mq.connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.config.RabbitMQ.Queue, // name
		false,                    // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		mq.config.RabbitMQ.Exchange, // exchange
		q.Name,                      // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	return nil
}
