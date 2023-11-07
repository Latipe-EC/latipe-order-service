package message

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"order-rest-api/config"
	"time"
)

type ProducerOrderMessage struct {
	channel *amqp.Channel
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

	return nil
}

func SendMessage(request interface{}, content ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(content) > 1 {
		content[0] = "default"
	}

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info]: Send message to queue %v", content[0])
	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.Exchange,   // exchange
		producer.cfg.RabbitMQ.RoutingKey, // routing key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return nil
}
