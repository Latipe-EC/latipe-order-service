package message

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"order-service-rest-api/config"
)

func NewConnectionToRabbitMQ(config *config.Config) *amqp.Connection {
	conn, err := amqp.Dial(config.RabbitMQ.Connection)
	failOnError(err, "Failed to connect to RabbitMQ")

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	return conn
}
