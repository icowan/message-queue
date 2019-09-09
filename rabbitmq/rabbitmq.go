package rabbitmq

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/icowan/message-queue/queue"
	"github.com/streadway/amqp"
)

type rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *RabbitMqConfig
}

type RabbitMqConfig struct {
	Url          string `json:"url,omitempty"`
	Exchange     string `json:"exchange,omitempty"`
	ExchangeType string `json:"exchange_type,omitempty"`
	RoutingKey   string `json:"routing_key,omitempty"`
}

func NewRabbitMq(config *RabbitMqConfig) (mq queue.MessageQueue, err error) {
	q := &rabbitmq{}
	if config == nil {
		config = &RabbitMqConfig{
			Url: "amqp://guest:guest@127.0.0.1:5672/test",
		}
	}
	q.config = config
	if err := q.connect(); err != nil {
		return mq, err
	}
	return q, nil
}

func (c *rabbitmq) connect() error {
	if c.conn == nil {
		conn, err := amqp.Dial(c.config.Url)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	return nil
}

func (c *rabbitmq) Publish(exchangeName string, exchangeType string, data []byte, contentType string) (err error) {
	if err = c.connect(); err != nil {
		return
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return
	}

	q, err := ch.QueueDeclare("", false, false, false, false, nil)
	if err != nil {
		return
	}

	if err = ch.QueueBind(q.Name, exchangeName, exchangeName, false, nil); err != nil {
		return
	}

	if err = ch.Publish(exchangeName, exchangeName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Transient,
		Priority:     0,
		ContentType:  contentType,
		Body:         data,
	}); err != nil {
		return
	}
	return
}

func (c *rabbitmq) PublishOnQueue(queueName string, data []byte, contentType string) (err error) {
	if err = c.connect(); err != nil {
		return
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return
	}
	defer func() {
		if err = ch.Close(); err != nil {
			fmt.Errorf("ch Close err %v", err)
		}
	}()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return
	}

	if err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Transient,
		Priority:     0,
		ContentType:  "application/json",
		Body:         data,
	}); err != nil {
		return
	}
	return
}

func (c *rabbitmq) Subscribe(exchangeName string, exchangeType string, consumerName string, f func(data string)) error {
	return nil
}

func (c *rabbitmq) SubscribeToQueue(queueName, consumerName string, f func(data string)) (err error) {
	if err = c.connect(); err != nil {
		return
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return
	}

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return
	}

	deliveries, err := ch.Consume(q.Name, consumerName, true, false, false, false, nil)
	if err != nil {
		return
	}

	if err = ch.QueueBind(queueName, c.config.RoutingKey, c.config.Exchange, false, nil); err != nil {
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range deliveries {
			s := BytesToString(&(d.Body))
			f(*s)
		}
	}()

	<-forever
	return
}

func (c *rabbitmq) Close() error {
	if c.conn != nil {
		defer func() {
			if err := c.conn.Close(); err != nil {
				fmt.Errorf("ch Close err %v", err)
			}
		}()
	}
	return errors.New("mq not connected.")
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
