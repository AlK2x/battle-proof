package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"jam/config"
	"jam/internal/domain/jam"
	"time"

	"github.com/segmentio/kafka-go"
)

const JamCreatedTopic = "jam.created"
const ParticipantAddedTopic = "participant.added"

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
	payload, topic, err := f.marshalPayload(event)
	if err != nil {
		return kafka.Message{}, err
	}
	return kafka.Message{
		Key:   []byte(getEventKey(event)),
		Value: payload,
		Topic: topic,
	}, nil
}

func (f *KafkaEventFactory) marshalPayload(event jam.DomainEvent) ([]byte, string, error) {
	switch e := event.(type) {
	case *jam.JamCreatedEvent:
		payload, err := json.Marshal(jamCreatedEvent{
			ID:         string(e.ID),
			Type:       string(e.Type),
			CreatedBy:  e.CreatedBy,
			OccurredAt: e.OccurredAt,
			Name:       e.Name,
			Location:   e.Location,
		})
		return payload, JamCreatedTopic, err
	case *jam.ParticipantAddedEvent:
		payload, err := json.Marshal(participantAddedEvent{
			ID:         string(e.ID),
			Type:       string(e.Type),
			CreatedBy:  e.CreatedBy,
			OccurredAt: e.OccurredAt,
			UserID:     e.UserID,
		})
		return payload, ParticipantAddedTopic, err
	}
	return []byte{}, "", fmt.Errorf("unknown DomainEvent type: %T", event)
}

func getEventKey(event jam.DomainEvent) string {
	switch e := event.(type) {
	case *jam.JamCreatedEvent:
		return string(e.ID)
	case *jam.ParticipantAddedEvent:
		return string(e.ID)
	default:
		return ""
	}
}

func NewKafkaEventProducer(config config.Config) *KafkaEventProducer {
	return &KafkaEventProducer{
		messageFactory: KafkaEventFactory{},
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(config.KafkaAddr),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

type KafkaEventProducer struct {
	messageFactory KafkaEventFactory
	writer         *kafka.Writer
}

func (p *KafkaEventProducer) Publish(ctx context.Context, event jam.DomainEvent) error {
	e, err := p.messageFactory.Create(event)
	if err != nil {
		return err
	}
	err = p.writer.WriteMessages(ctx, e)
	return err
}

func (p *KafkaEventProducer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
