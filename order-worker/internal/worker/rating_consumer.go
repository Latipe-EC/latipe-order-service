package worker

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"order-worker/config"
	"order-worker/internal/app/orders"
	messageDTO "order-worker/internal/domain/dto"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerRatingMessage struct {
	config       *config.Config
	orderUsecase orders.Usecase
}

func NewConsumerRatingMessage(config *config.Config, orderService orders.Usecase) *ConsumerRatingMessage {
	return &ConsumerRatingMessage{
		config:       config,
		orderUsecase: orderService,
	}
}

func (mq ConsumerRatingMessage) ListenRatingEventQueue(wg *sync.WaitGroup) {
	conn, err := amqp.Dial(mq.config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Info("Comsumer has been connected")

	channel, err := conn.Channel()
	defer channel.Close()
	defer conn.Close()

	// Khai báo một Exchange loại "direct"
	err = channel.ExchangeDeclare(
		mq.config.RabbitMQ.RatingEvent.Exchange, // Tên Exchange
		"topic",
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
		mq.config.RabbitMQ.RatingEvent.Queue,
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
		mq.config.RabbitMQ.RatingEvent.Queue,
		mq.config.RabbitMQ.RatingEvent.RoutingKey,
		mq.config.RabbitMQ.RatingEvent.Exchange,
		false,
		nil)
	if err != nil {
		log.Fatalf("cannot bind exchange: %v", err)
	}

	// declaring consumer with its properties over channel opened
	msgs, err := channel.Consume(
		mq.config.RabbitMQ.RatingEvent.Queue, // queue
		mq.config.RabbitMQ.ConsumerName,      // consumer
		true,                                 // auto ack
		false,                                // exclusive
		false,                                // no local
		false,                                // no wait
		nil,                                  //args
	)
	if err != nil {
		panic(err)
	}

	defer wg.Done()
	// handle consumed messages from queue
	for msg := range msgs {
		log.Infof(" received order message from: %s", msg.RoutingKey)

		if err := mq.handleMessage(msg); err != nil {
			log.Errorf(" updating the rating was failed cause %s", err)
		}

	}

	log.Info("message queue has started")
	log.Info("waiting for messages...")

}

func (mq ConsumerRatingMessage) handleMessage(msg amqp.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := messageDTO.RatingMessage{}

	if err := json.Unmarshal(msg.Body, &message); err != nil {
		log.Errorf("parse rating message was failed cause: %s", err)
		return err
	}

	err := mq.orderUsecase.UpdateRatingItem(ctx, &message)
	if err != nil {
		return err
	}

	log.Infof("the rating was updated successfully: %s", msg.RoutingKey)
	return nil
}
