package main

import (
	"fmt"
	"log"
	server "order-worker/internal"
	"sync"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[%s] Application has closed")

	serv, err := server.New()
	if err != nil {
		fmt.Printf("%s", err)
	}

	//order handle worker
	var wg sync.WaitGroup
	wg.Add(1)
	go serv.ConsumerOrderMessage().ListenOrderEventQueue(&wg)

	wg.Add(1)
	go serv.ConsumerRatingMessage().ListenRatingEventQueue(&wg)

	wg.Add(1)
	go serv.OrderCompleteCJ().CheckOrderFinishShippingStatus(&wg)

	wg.Wait()
}
