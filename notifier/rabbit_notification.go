package notifier

type RabbitNotification struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Body       []byte
}
