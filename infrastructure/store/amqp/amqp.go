package amqp

import (
	"github.com/streadway/amqp"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository"
)

type store struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(connString string) repository.AMQPStore {
	connection, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare("change_bid", false, false, false, false, nil)
	if err != nil {
		return err
	}
	return &store{connection, channel}
}

func (s *store) Publish(msg string) error {
	//msg := amqp.Publishing{
	//
	//}
	//return s.channel.Publish("", "change_bid", false, false, msg)
	return nil
}

func (s *store) Shutdown() error {
	return s.connection.Close()
}
