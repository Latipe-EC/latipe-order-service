package message

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"

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
	log.Info("producer has been connected")

	producer.cfg = config
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return err
	}
	producer.channel = ch

	return nil
}

func SendOrderMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		producer.cfg.RabbitMQ.OrderEvent.Exchange,
		producer.cfg.RabbitMQ.OrderEvent.RoutingKey)

	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.OrderEvent.Exchange,   // exchange
		producer.cfg.RabbitMQ.OrderEvent.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return nil
}

func SendCancelingOrderMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Infof("Send message to queue %v - %v",
		producer.cfg.RabbitMQ.PaymentEvent.Exchange,
		producer.cfg.RabbitMQ.PaymentEvent.RoutingKey)

	err = producer.channel.PublishWithContext(ctx,
		producer.cfg.RabbitMQ.PaymentEvent.Exchange,   // exchange
		producer.cfg.RabbitMQ.PaymentEvent.RoutingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return nil
}
