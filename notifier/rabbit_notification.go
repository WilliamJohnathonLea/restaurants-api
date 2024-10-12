package notifier

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitNotification struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Headers    amqp.Table
	Body       []byte
}
