package message

import (
	"context"
	"log"
	"order-worker/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderWorkerProducer struct {
	channel *amqp.Channel
	cfg     *config.Config
}

var producer OrderWorkerProducer

func InitWorkerProducer(config *config.Config) error {
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

	failOnError(err, "Failed to declare a queue")

	return nil
}

func SendEmailMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info]: Send message %v", request)
	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.EmailEvent.Exchange,   // exchange
		producer.cfg.RabbitMQ.EmailEvent.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a email message")

	return nil
}

func SendCartServiceMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info]: Send message %v", request)
	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.CartEvent.Exchange,   // exchange
		producer.cfg.RabbitMQ.CartEvent.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a email message")

	return nil
}
