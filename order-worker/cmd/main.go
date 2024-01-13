package main

import (
	"fmt"
	"log"
	server "order-worker/internal"
	"sync"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	//order handle worker
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	wg.Add(1)
	go serv.ConsumerOrderMessage().ListenOrderEventQueue(&wg)

	wg.Add(1)
	go serv.ConsumerRatingMessage().ListenRatingEventQueue(&wg)

	wg.Add(1)
	go serv.OrderCompleteCJ().CheckOrderFinishShippingStatus(&wg)

	wg.Wait()
}
