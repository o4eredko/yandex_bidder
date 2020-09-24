package amqp

import (
	"github.com/streadway/amqp"

	amqpRepo "gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
)

type store struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(config *config.AMQP) amqpRepo.Store {
	connection, err := amqp.Dial(config.DSN())
	if err != nil {
		panic(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	return &store{
		conn:    connection,
		channel: channel,
	}
}

func (s *store) Publish(publishConfig *config.PublishConfig, message *amqp.Publishing) error {
	if err := s.channel.ExchangeDeclare(
		publishConfig.Exchange.Name,
		publishConfig.Exchange.Type,
		publishConfig.Exchange.Durable,
		publishConfig.Exchange.AutoDelete,
		publishConfig.Exchange.Internal,
		publishConfig.Exchange.NoWait,
		publishConfig.Exchange.Args,
	); err != nil {
		return err
	}

	err := s.channel.Publish( // Publishes a message onto the queue.
		publishConfig.Exchange.Name,
		publishConfig.Exchange.RoutingKey,
		publishConfig.Mandatory,
		publishConfig.Immediate,
		*message,
	)
	return err
}

func (s *store) Subscribe(consumeConfig *config.ConsumeConfig, handler func(amqp.Delivery)) error {
	if err := s.channel.ExchangeDeclare(
		consumeConfig.Exchange.Name,
		consumeConfig.Exchange.Type,
		consumeConfig.Exchange.Durable,
		consumeConfig.Exchange.AutoDelete,
		consumeConfig.Exchange.Internal,
		consumeConfig.Exchange.NoWait,
		consumeConfig.Exchange.Args,
	); err != nil {
		return err
	}

	queue, err := s.channel.QueueDeclare(
		consumeConfig.Exchange.Queue.Name,
		consumeConfig.Exchange.Queue.Durable,
		consumeConfig.Exchange.Queue.AutoDelete,
		consumeConfig.Exchange.Queue.Exclusive,
		consumeConfig.Exchange.Queue.NoWait,
		consumeConfig.Exchange.Queue.Args,
	)
	if err != nil {
		return err
	}

	if err = s.channel.QueueBind(
		consumeConfig.Exchange.Queue.Name,
		consumeConfig.Exchange.RoutingKey,
		consumeConfig.Exchange.Name,
		consumeConfig.Exchange.Queue.NoWait,
		consumeConfig.Exchange.Queue.Args,
	); err != nil {
		return err
	}

	messages, err := s.channel.Consume(
		queue.Name,
		consumeConfig.Name,
		consumeConfig.AutoAck,
		consumeConfig.Exclusive,
		consumeConfig.NoLocal,
		consumeConfig.NoWait,
		consumeConfig.Args,
	)
	if err != nil {
		return err
	}

	go consumeLoop(messages, handler)
	return nil
}

func (s *store) Shutdown() error {
	s.channel.Close()
	err := s.conn.Close()
	return err
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}
