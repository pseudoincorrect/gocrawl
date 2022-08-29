package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gocrawl/pkg/mq"
	"gocrawl/pkg/theprinter"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type App struct {
	*mux.Router
	*mq.MsgQueue
}

var app App

func main() {
	theprinter.SayWhoYouAre("SERVER")
	app = App{}
	app.initializeRouter()
	app.initializeMQ()
	app.startRouter()
	forever := make(chan interface{})
	<-forever
}

// startRouter start the HTTP router
func (a *App) startRouter() {
	http.ListenAndServe(":8080", a.Router)
}

// initializeRouter initialize the HTTP router and its handlers
func (a *App) initializeRouter() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/crawl", getCrawAnURLHandler(a)).
		Methods("POST")
}

// initializeMQ initialize the connection to RabbitMQ
func (a *App) initializeMQ() {
	msgQ, err := mq.Init("amqp://user:bitnami@rabbitmq:5672/")
	if err != nil {
		log.Fatalln("error initializing rabbitMQ:", err)
	}
	a.MsgQueue = msgQ
	log.Println("Connected to RabbitMQ")
}

type crawlURLRequest struct {
	URL string `json:"url"`
}

type httpHandler func(w http.ResponseWriter, r *http.Request)

// getCrawAnURLHandler return a HTTP handler for reauests to crawl
// an URL. Send a task to the MQ
func getCrawAnURLHandler(a *App) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		body := crawlURLRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Body :", body)
		if _, err = url.ParseRequestURI(body.URL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = a.addToQueue(body.URL)
		if err != nil {
			log.Println("could not add to queue err:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type crawlMsg struct {
	Url      string `json:"url"`
	Depth    int    `json:"depth"`
	MaxDepth int    `json:"max-depth"`
}

// addToQueue send a url to the MQ to be used by worker thereafter
func (a *App) addToQueue(url string) error {
	q, err := a.Channel.QueueDeclare(
		"crawl_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Println("QueueDeclare err:", err)
		return err
	}
	msg := crawlMsg{
		Url:      url,
		Depth:    0,
		MaxDepth: 1,
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Println("Json marshal err:", err)
		return err
	}
	ctx := context.Background()
	err = a.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msgJson,
		})
	fmt.Printf("msgJson = %v\n", msgJson)
	if err != nil {
		fmt.Println("publish err:", err)
		return err
	}
	return nil
}
