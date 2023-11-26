package publisher

import (
	"context"
	"log"
	"order-worker/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageProducer struct {
	channel *amqp.Channel
	cfg     *config.Config
}

func InitWorkerProducer(config *config.Config) *MessageProducer {
	conn, err := amqp.Dial(config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Printf("[%s] Producer has been connected", "INFO")

	producer := MessageProducer{
		channel: nil,
		cfg:     config,
	}

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
		return nil
	}

	producer.channel = ch

	return &producer
}

func (p MessageProducer) SendEmailMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	err = p.channel.ExchangeDeclare(
		p.cfg.RabbitMQ.EmailEvent.Exchange, // name
		"topic",                            // type
		true,                               // durable
		false,                              // auto-deleted
		false,                              // internal
		false,                              // no-wait
		nil,                                // arguments
	)

	log.Printf("[Info]: Send message %v", request)
	err = p.channel.PublishWithContext(ctx,
		p.cfg.RabbitMQ.EmailEvent.Exchange,   // exchange
		p.cfg.RabbitMQ.EmailEvent.RoutingKey, // routing key
		false,                                // mandatory
		false,                                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a email message")

	return nil
}

func (p MessageProducer) SendCartServiceMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info] [%v-%v]: Send message %v",
		p.cfg.RabbitMQ.CartEvent.Exchange,
		p.cfg.RabbitMQ.CartEvent.RoutingKey,
		request)

	err = p.channel.PublishWithContext(ctx,
		p.cfg.RabbitMQ.CartEvent.Exchange,   // exchange
		p.cfg.RabbitMQ.CartEvent.RoutingKey, // routing key
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a email message")

	return nil
}

func (p MessageProducer) SendBillingServiceMessage(request interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := ParseOrderToMessage(&request)
	if err != nil {
		return err
	}

	log.Printf("[Info] [%v-%v]: Send message %v",
		p.cfg.RabbitMQ.BillingEvent.Exchange,
		p.cfg.RabbitMQ.BillingEvent.RoutingKey,
		request)

	err = p.channel.PublishWithContext(ctx,
		p.cfg.RabbitMQ.BillingEvent.Exchange,   // exchange
		p.cfg.RabbitMQ.BillingEvent.RoutingKey, // routing key
		false,                                  // mandatory
		false,                                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a email message")

	return nil
}
