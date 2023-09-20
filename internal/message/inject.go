package message

import "github.com/google/wire"

var Set = wire.NewSet(
	NewConnectionToRabbitMQ,
	NewProducerOrderMessage,
	NewConsumerOrderMessage,
)
