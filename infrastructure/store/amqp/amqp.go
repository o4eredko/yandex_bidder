package amqp

import (
	"github.com/streadway/amqp"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
)

type (
	Store struct {
		connection   *amqp.Connection
		exchangeName string
	}
)

func New(config *config.AMQP) *Store {
	connection, err := amqp.Dial(config.DSN())
	if err != nil {
		panic(err)
	}

	return &Store{
		connection:   connection,
		exchangeName: "change_bid",
	}
}

func (s *Store) Publish(msg []byte) error {
	channel, err := s.connection.Channel()
	if err != nil {
		return err
	}
	if err := channel.ExchangeDeclare(
		s.exchangeName,
		"x-delayed-message",
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": amqp.ExchangeTopic,
		},
	); err != nil {
		return err
	}

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	}

	return channel.Publish(s.exchangeName, "bid.update", false, false, publishing)
}

func (s *Store) Shutdown() error {
	return s.connection.Close()
}
