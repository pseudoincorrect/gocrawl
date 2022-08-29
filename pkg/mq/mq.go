package mq

import (
	"fmt"
	"gocrawl/pkg/apperror"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MsgQueue struct {
	*amqp.Connection
	*amqp.Channel
}

func Init(address string) (*MsgQueue, error) {
	fmt.Println("rabitmq")
	conn, err := amqp.Dial(address)
	apperror.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	apperror.FailOnError(err, "Failed to open a channel")

	return &MsgQueue{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (m *MsgQueue) CloseChannel() {
	m.Channel.Close()
}

func (m *MsgQueue) CloseConnection() {
	m.Connection.Close()
}
