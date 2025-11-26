package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/app"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/storage"
)

// Event represents the structure of an event.
type Event struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type Producer struct {
	app      *app.App
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	key      string
}

// NewProducer initializes the RabbitMQ producer and associates it with the app instance.
func NewProducer(a *app.App, uri, exchange, key string) (*Producer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Producer{
		app:      a,
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		key:      key,
	}, nil
}

func (p *Producer) Start(quit <-chan struct{}) {
	// Create a new cron schedulers
	c := cron.New() // supports seconds like "every 10s" intervals

	// 1. Schedule hourly job (configurable interval would be better)
	_, err := c.AddFunc("@every 1h", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("[Producer] Checking for events to publish...")

		if err := p.ListEventsDay(ctx); err != nil {
			log.Printf("[Producer] failed to publish daily events: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("[Producer] Failed to schedule cron: %v", err)
	}

	// 2. Clean old events once a day at 03:00
	_, err = c.AddFunc("0 3 * * 0", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		log.Println("[Producer] Running daily cleanup to remove 1 year old events...")

		if err := p.CleanOldEvents(ctx); err != nil {
			log.Printf("[Producer] failed to clean events: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("[Producer] Failed to schedule cleanup cron: %v", err)
	}
	c.Start()
	log.Println("[Producer] Cron started: publishing every 1h")

	// Block until quit signal is received
	<-quit

	log.Println("[Producer] Shutdown signal received, stopping cron...")
	c.Stop()
	log.Println("[Producer] Cron stopped gracefully.")
}

// ListEventsDay generates and publishes event data for the day.
func (p *Producer) ListEventsDay(ctx context.Context) error {
	events, err := p.app.ListEvents(ctx, storage.PeriodDay)
	if err != nil {
		log.Printf("failed to list day events: %v", err)
		return fmt.Errorf("failed to list day events: %w", err)
	}

	log.Printf("[Producer] found %d events", len(events))
	for _, event := range events {
		if event.Start.Hour() != time.Now().Hour() {
			log.Printf("[Producer] skipped old event: Title=%s Start=%s", event.Title, event.Start)
			continue
		}

		msg, err := json.Marshal(event)
		if err != nil {
			log.Printf("[Producer] failed to serialize event: %v", err)
			continue
		}

		if err := p.Publish(msg); err != nil {
			log.Printf("[Producer] failed to publish: %v", err)
			continue
		}

		log.Printf("[Producer] sent: %s", msg)
	}
	return nil
}

// Publish sends a raw message (JSON) to RabbitMQ.
func (p *Producer) Publish(body []byte) error {
	return p.channel.Publish(
		p.exchange,
		p.key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// CleanOldEvents removes events older than one year.
func (p *Producer) CleanOldEvents(ctx context.Context) error {
	events, err := p.app.ListEvents(ctx, storage.PeriodAll)
	if err != nil {
		log.Printf("failed to get all events: %v", err)
		return fmt.Errorf("failed to get all events: %w", err)
	}

	log.Printf("[Producer] found %d events", len(events))
	for _, event := range events {
		// Delete Event if it is a year old.
		if event.Start.Before(time.Now().AddDate(-1, 0, 0)) {
			log.Printf("[Producer] deleting old event: Title=%s Start=%s", event.Title, event.Start)
			p.app.DeleteEvent(ctx, event.ID)
		}
	}
	return nil
}

func (p *Producer) Shutdown() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
