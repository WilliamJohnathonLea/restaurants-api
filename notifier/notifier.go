package notifier

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitOpts func(*RabbitNotifer)

type RabbitNotifer struct {
	amqpUrl       string
	conn          *amqp.Connection
	channel       *amqp.Channel
	confirmations chan amqp.Confirmation
}

func (rn *RabbitNotifer) Close() error {
	chanErr := rn.channel.Close()
	connErr := rn.conn.Close()

	return errors.Join(chanErr, connErr)
}

func (rn *RabbitNotifer) Notify(notification RabbitNotification) error {
	// The exchange, queue and bind are assumed to be in place
	err := rn.channel.Publish(
		notification.Exchange,
		notification.RoutingKey,
		notification.Mandatory,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        notification.Body,
		},
	)

	// Return errors not related to publishing, e.g. connectivity
	if err != nil {
		return err
	}

	conf := <-rn.confirmations
	if conf.Ack {
		return nil
	} else {
		return errors.New("message was not delivered")
	}

	// TODO Check for message returns
}

func WithURL(url string) RabbitOpts {
	return func(rn *RabbitNotifer) {
		rn.amqpUrl = url
	}
}

func NewRabbitNotifier(opts ...RabbitOpts) (*RabbitNotifer, error) {
	rn := &RabbitNotifer{}

	for _, opt := range opts {
		opt(rn)
	}

	rabbitConn, err := amqp.Dial(rn.amqpUrl)
	if err != nil {
		return nil, err
	}

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		return nil, err
	}

	rn.conn = rabbitConn
	rn.channel = rabbitCh

	// Using confirm mode to check publishings to the server
	err = rn.channel.Confirm(false)
	if err != nil {
		rn.Close()
		return nil, err
	}

	// Channel for publish notifications
	rn.confirmations = rn.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	return rn, nil
}
