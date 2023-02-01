package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// step-1 : try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// step-2 : start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// step-3 : create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// step-4 : watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	// rabbitmq can be slow start
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

	// don't continue until rabbit is ready
	for {
		// 'rabbitmq' is the service name of the rabbitmq container in docker-compose file
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready ...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// Each time the connection fails, increase the delay
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing off ...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
