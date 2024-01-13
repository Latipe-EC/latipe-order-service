package main

import (
	"fmt"
	"log"
	server "order-rest-api/internal"
	"order-rest-api/internal/message"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	//init message queue
	if err := message.InitProducerMessage(serv.Config()); err != nil {
		fmt.Printf("%s", err)
	}

	if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
		fmt.Printf("%s", err)
	}
}
