package main

import (
	"fmt"
	"order-service-rest-api/config"
	"order-service-rest-api/internal/infrastructure/persistence/db"
)

func main() {
	fmt.Println("Init application")
	configuration, err := config.NewConfig()
	if err != nil {
		fmt.Printf("\nERROR:%v", err.Error())
	}
	gorm := db.NewMySQLConnection(configuration)
	if gorm != nil {
		fmt.Println(gorm.Name())
	}
}
