package amqp

import (
	"github.com/streadway/amqp"
)

type (
	Store struct {
		connection *amqp.Connection
		Channel    *amqp.Channel
	}
)

func New(connString string) *Store {
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
		panic(err)
	}
	return &Store{
		connection: connection,
		Channel:    channel,
	}
}

func (s *Store) Shutdown() error {
	return s.connection.Close()
}
