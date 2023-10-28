package main

import (
	"fmt"
	"log"
	server "order-worker/internal"
	"order-worker/internal/message"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[%s] Application has closed")

	serv, err := server.New()
	if err != nil {
		fmt.Printf("%s", err)
	}

	//order handle worker
	go func() {
		serv.Consumer().ListenMessageQueue()
	}()

	//init message queue
	if err := message.InitWorkerProducer(serv.Config()); err != nil {
		fmt.Printf("%s", err)
	}

	if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
		fmt.Printf("%s", err)
	}
}
