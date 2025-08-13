package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/akkahshh24/movieapp/rating/pkg/model"
	"github.com/segmentio/kafka-go"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	reader *kafka.Reader
}

// NewIngester creates a new Kafka ingester.
func NewIngester(addr string, groupID string, topic string) (*Ingester, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{addr},
		GroupID:     groupID,
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
		MaxWait:     time.Second,
	})
	return &Ingester{reader: reader}, nil
}

// Ingest starts reading messages from Kafka and sends them over a channel.
func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	fmt.Println("Starting Kafka ingester")

	ch := make(chan model.RatingEvent, 1)
	go func() {
		defer close(ch)
		defer i.reader.Close()

		for {
			m, err := i.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Println("Kafka read error:", err)
				continue
			}

			var event model.RatingEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				fmt.Println("Unmarshal error:", err)
				continue
			}
			ch <- event
		}
	}()

	return ch, nil
}
