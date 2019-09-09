package queue

type MessageQueue interface {
	Close() error
	Publish(exchangeName string, exchangeType string, data []byte, contentType string) (err error)
	Subscribe(exchangeName string, exchangeType string, consumerName string, f func(data string)) error
	PublishOnQueue(queueName string, data []byte, contentType string) (err error)
	SubscribeToQueue(queueName, consumerName string, f func(data string)) (err error)
}
