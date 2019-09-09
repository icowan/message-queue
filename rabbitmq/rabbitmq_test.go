package rabbitmq

import (
	"fmt"
	"testing"
	"time"
)

func TestRabbitmq_PublishOnQueue(t *testing.T) {

	q, err := NewRabbitMq(&RabbitMqConfig{
		Url:      "amqp://kplcloud:helloworld@rabbitmq:5672/kplcloud",
		Exchange: "kplcloud",
	})
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err = q.Close(); err != nil {
			fmt.Errorf("ch Close err %v", err)
		}
	}()

	if err = q.PublishOnQueue("test", []byte("111222"), ""); err != nil {
		t.Error(err)
	}

	go func() {
		if err = q.SubscribeToQueue("test", "", func(data string) {
			t.Log(data)
			fmt.Println(data)
		}); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(20 * time.Second)
}
