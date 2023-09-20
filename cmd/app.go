package main

import (
	"fmt"
	server "order-service-rest-api/internal"
)

func main() {
	fmt.Println("Init application")
	serv, err := server.New()
	if err != nil {
		fmt.Printf("%s", err)
	}
	if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
		fmt.Printf("%s", err)
	}
}
