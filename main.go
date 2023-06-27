package main

import (
	"fmt"
	"log"

	database "github.com/ponyjackal/order-book/Database"
	env "github.com/ponyjackal/order-book/Helper"
	kafka "github.com/ponyjackal/order-book/Kafka"
	orderModel "github.com/ponyjackal/order-book/Models"
	router "github.com/ponyjackal/order-book/Router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	_, err := database.Connect(database.DBconnectionConfig{
		Url:      env.Getenv("DB_URL", "localhost"),
		Port:     env.Getenv("DB_PORT", "3306"),
		Username: env.Getenv("DB_USERNAME", "root"),
		Password: env.Getenv("DB_PASSWORD", ""),
		DBname:   env.Getenv("DB_NAME", ""),
	})

	if err != nil {
		log.Fatalln(err)
	}

	go startConsumeOrder()
	
	ginEngine := gin.Default()
	router.DefineRoutes(ginEngine)
	ginEngine.Run()
}

func startConsumeOrder() {
	kafkaConnection := kafka.Connection()
	defer kafkaConnection.Close()

	fmt.Println("Listen for kafka data....")
	for {
		var order orderModel.Order
		var err error
		order, err = kafka.ConsumeOrder(kafkaConnection)
		if err != nil {
			fmt.Println(err)
			continue
		}
		order.Insert()
		fmt.Println("Consome order...")
		fmt.Println(order)
	}
}
