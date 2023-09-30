package message

import (
	"context"
	"log"
	"order-worker/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProducerOrderMessage struct {
	channel *amqp.Channel
	queue   *amqp.Queue
	cfg     *config.Config
}

var producer ProducerOrderMessage

func InitProducerMessage(config *config.Config) error {
	conn, err := amqp.Dial(config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] Producer has been connected", "INFO")

	producer.cfg = config
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}
	producer.channel = ch

	q, err := producer.channel.QueueDeclare(
		config.RabbitMQ.Queue, // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	producer.queue = &q
	failOnError(err, "Failed to declare a queue")

	return nil
}

func SendMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info]: Send message %v", request)
	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.Exchange, // exchange
		producer.queue.Name,            // routing key
		false,                          // mandatory
		false,                          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return nil
}
