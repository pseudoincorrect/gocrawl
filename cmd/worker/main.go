package main

import (
	"encoding/json"
	"fmt"
	"gocrawl/pkg/mq"
	"gocrawl/pkg/theprinter"
	"log"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	*mux.Router
	*mq.MsgQueue
}

var app App

func main() {
	app = App{}
	app.initializeMQ()
	go sayWhoYouAre()
	forever := make(chan interface{})
	<-forever
}

func sayWhoYouAre() {
	for {
		theprinter.SayWhoYouAre("SUPER WORKER")
		time.Sleep(5 * time.Second)
	}
}

type CrawlMsg struct {
	Url      string `json:"url"`
	Depth    int    `json:"depth"`
	MaxDepth int    `json:"max-depth"`
}

// initializeMQ initialise the connection to RabbitMQ
func (a *App) initializeMQ() {
	msgQ, err := mq.Init("amqp://user:bitnami@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("error initializing rabbitMQ:", err)
	}
	a.MsgQueue = msgQ

	q, err := a.Channel.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("error declaring queue")
	}

	err = a.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal("Failed to set QoS")
	}

	msgs, err := a.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			crawlMsg := CrawlMsg{}
			err := json.Unmarshal(d.Body, &crawlMsg)
			if err != nil {
				log.Println("error unmarchalling msg err:", err)
				return
			}
			indent, _ := json.MarshalIndent(crawlMsg, "", "  ")
			fmt.Println("crawl msg :", indent)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// func main() {
// 	url := flag.String("url", "example.com", "website url")
// 	flag.Parse()
// 	fmt.Println("url = ", *url)
// 	body := fetch.Fetch(*url)
// 	reader := strings.NewReader(body)
// 	links, err := parser.Parse(reader)
// 	if err != nil {
// 		log.Fatal("parser error")
// 	}
// 	fmt.Println(links)
// }
