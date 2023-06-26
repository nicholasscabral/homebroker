package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/nicholasscabral/homebroker/go/internal/infra/kafka"
	"github.com/nicholasscabral/homebroker/go/internal/market/dto"
	"github.com/nicholasscabral/homebroker/go/internal/market/entity"
	"github.com/nicholasscabral/homebroker/go/internal/market/transformer"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChannel := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "asdf",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewProducer(configMap)
	consumer := kafka.NewConsumer(configMap, []string{"orders"}) // orders = input

	go consumer.Consume(kafkaMsgChannel) // new thread to read information from channel

	// receive on kafka channel => process the input => publish on kafka
	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // new thread to listen for the trades

	go func() {
		for msg := range kafkaMsgChannel {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
