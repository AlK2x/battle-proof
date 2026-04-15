package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"jam/config"
	"jam/internal/domain/jam"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const Topic = "jams"

type jamCreatedEvent struct {
	ID         string    `json:"event_id"`
	Type       string    `json:"event_type"`
	CreatedBy  string    `json:"created_by"`
	OccurredAt time.Time `json:"occured_at"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
}

type participantAddedEvent struct {
	ID         string    `json:"event_id"`
	Type       string    `json:"event_type"`
	CreatedBy  string    `json:"created_by"`
	OccurredAt time.Time `json:"occured_at"`
	UserID     string    `json:"user_id"`
}

type KafkaEventFactory struct {
}

func (f *KafkaEventFactory) Create(event jam.DomainEvent) (kafka.Message, error) {
	payload, err := f.marshalPayload(event)
	if err != nil {
		return kafka.Message{}, err
	}
	return kafka.Message{
		Value: payload,
	}, nil
}

func (f KafkaEventFactory) marshalPayload(event jam.DomainEvent) ([]byte, error) {
	switch e := event.(type) {
	case *jam.JamCreatedEvent:
		return json.Marshal(jamCreatedEvent{
			ID:         string(e.ID),
			Type:       string(e.Type),
			CreatedBy:  e.CreatedBy,
			OccurredAt: e.OccurredAt,
			Name:       e.Name,
			Location:   e.Location,
		})
	case *jam.ParticipantAddedEvent:
		return json.Marshal(participantAddedEvent{
			ID:         string(e.ID),
			Type:       string(e.Type),
			CreatedBy:  e.CreatedBy,
			OccurredAt: e.OccurredAt,
			UserID:     e.UserID,
		})
	}
	return []byte{}, errors.New("unknown DomainEventType")
}

var kafkaPublisher *KafkaEventProducer
var once sync.Once

func GetKafkaEventProducerInstance(ctx context.Context, config config.Config) *KafkaEventProducer {
	if kafkaPublisher != nil {
		return kafkaPublisher
	}

	once.Do(func() {
		kafkaPublisher = &KafkaEventProducer{
			context:        ctx,
			messageFactory: KafkaEventFactory{},
			writer: &kafka.Writer{
				Addr:     kafka.TCP(config.KafkaAddr),
				Topic:    Topic,
				Balancer: &kafka.LeastBytes{},
			},
		}
	})
	return kafkaPublisher
}

type KafkaEventProducer struct {
	context        context.Context
	messageFactory KafkaEventFactory
	writer         *kafka.Writer
}

func (p *KafkaEventProducer) Publish(event jam.DomainEvent) error {
	e, err := p.messageFactory.Create(event)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(p.context, e)
	return err
}
