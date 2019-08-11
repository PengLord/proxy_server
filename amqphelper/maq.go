package amqphelper

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Config struct {
	AmqpUrl    string
	WorkerName string
	FromQueue  FromQueue
}

type FromQueue struct {
	Name          string
	ExchangeName  string
	Key           string
	PrefetchCount int
	Priority      int
}
type TaskQueueDispatcher struct {
	Config          *Config
	connection      *amqp.Connection
	consumerChannel *amqp.Channel
	pushChannel     *amqp.Channel
}

func New(config *Config) TaskQueueDispatcher {
	//1. 初始化操作
	worker := TaskQueueDispatcher{
		Config: config,
	}

	conn, err := amqp.Dial(worker.Config.AmqpUrl)
	if err != nil {
		log.Fatal(err)
	}
	consumerChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	pushChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	worker.connection = conn
	worker.consumerChannel = consumerChannel
	worker.pushChannel = pushChannel
	if worker.Config.FromQueue.PrefetchCount == 0 {
		worker.Config.FromQueue.PrefetchCount = 10
	}
	if worker.Config.FromQueue.Priority == 0 {
		worker.Config.FromQueue.Priority = 10
	}
	return worker
}
func (lwd *TaskQueueDispatcher) declare() <-chan amqp.Delivery {
	//1. 获取任务
	if lwd.Config.FromQueue.ExchangeName != "" {
		err := lwd.consumerChannel.ExchangeDeclare(
			lwd.Config.FromQueue.ExchangeName, // name
			"topic",                           // type
			true,                              // durable
			false,                             // auto-deleted
			false,                             // internal
			false,                             // no-wait
			nil,                               // arguments
		)
		failOnError(err, "Failed to declare an exchange")
	}
	priority := lwd.Config.FromQueue.Priority
	args := make(amqp.Table)
	args["x-max-priority"] = int16(priority)
	args["x-queue-mode"] = "lazy"
	q, err := lwd.consumerChannel.QueueDeclare(
		lwd.Config.FromQueue.Name, // name
		true,                      // durable
		false,                     // delete when unused
		false,                     // exclusive
		false,                     // no-wait
		args,                      // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = lwd.consumerChannel.Qos(
		lwd.Config.FromQueue.PrefetchCount, // prefetch count
		0,                                  // prefetch size
		false,                              // global
	)
	failOnError(err, "Failed to declare a Qos")
	if lwd.Config.FromQueue.ExchangeName != "" {
		err = lwd.consumerChannel.QueueBind(
			q.Name, // queue name
			lwd.Config.FromQueue.Key,
			lwd.Config.FromQueue.ExchangeName,
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}
	msgs, err := lwd.consumerChannel.Consume(
		lwd.Config.FromQueue.Name, // queue
		"",                        // consumer
		false,                     // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
	failOnError(err, "Failed to register a consumer")
	return msgs
}
func (lwd *TaskQueueDispatcher) StartConsume(runable func(buff []byte, delivery amqp.Delivery)) {
	msgs := lwd.declare()
	for d := range msgs {
		runable(d.Body, d)
	}
}
func (lwd *TaskQueueDispatcher) StartAsyncConsume(runable func(buff []byte, delivery amqp.Delivery)) {
	msgs := lwd.declare()
	//deliveryChan := make(chan amqp.Delivery, 30)
	//go func(deliveryChan chan amqp.Delivery) {
	//	for {
	//		select {
	//		case d := <-deliveryChan:
	//			go runable(d.Body, d)
	//		}
	//	}
	//}(deliveryChan)
	for d := range msgs {
		go runable(d.Body, d)
	}

}
func (lwd *TaskQueueDispatcher) Push(exchange, key string, priority uint8, body []byte) error {
	msg := amqp.Publishing{
		Priority:     priority,
		DeliveryMode: amqp.Persistent,
		Body:         body,
		ContentType:  "application/json",
	}
	return lwd.pushChannel.Publish(exchange, key, false, false, msg)

}
func (lwd *TaskQueueDispatcher) PushWithObject(exchange, key string, priority uint8, object interface{}) error {
	body, err := json.Marshal(object)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		Priority:     priority,
		DeliveryMode: amqp.Persistent,
		Body:         body,
		ContentType:  "application/json",
	}
	return lwd.pushChannel.Publish(exchange, key, false, false, msg)

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


