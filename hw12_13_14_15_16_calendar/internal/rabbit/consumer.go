package rabbit

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	tag     string
	msgs    <-chan amqp.Delivery
}

// NewConsumer sets up a RabbitMQ consumer.
func NewConsumer(uri, queue, tag string) (*Consumer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name,
		tag,
		false, // auto-ack = false
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   q.Name,
		tag:     tag,
		msgs:    msgs,
	}, nil
}

func (c *Consumer) Start(quit <-chan struct{}) {
	for {
		select {
		case msg := <-c.msgs:
			log.Printf("received: %s", msg.Body)
			msg.Ack(false)
		case <-quit:
			log.Println("consumer shutting down...")
			return
		}
	}
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
