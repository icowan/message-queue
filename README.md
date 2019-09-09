# message queue

## RabbitMq

```go
q, err := NewRabbitMq(&RabbitMqConfig{
    Url:      "amqp://kplcloud:helloworld@rabbitmq:5672/kplcloud",
    Exchange: "kplcloud",
})
if err != nil {
    t.Error(err)
}
defer q.Close()

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
```