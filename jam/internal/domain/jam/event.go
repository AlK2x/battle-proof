package jam

import (
	"context"
	"time"
)

type EventType string

const (
	JamCreated       EventType = "jam.created"
	ParticipantAdded EventType = "participant.added"
)

type EventID string

type DomainEvent interface {
	GetType() EventType
	GetOccurredAt() time.Time
}

type Event struct {
	ID         EventID
	Type       EventType
	CreatedBy  string
	OccurredAt time.Time
}

func (e *Event) GetType() EventType {
	return e.Type
}

func (e *Event) GetOccurredAt() time.Time {
	return e.OccurredAt
}

type JamCreatedEvent struct {
	Event
	Name     string
	Location string
	Date     time.Time
}

func (jce JamCreatedEvent) GetType() EventType {
	return jce.Event.Type
}

func (jce JamCreatedEvent) GetOccurredAt() time.Time {
	return jce.Event.OccurredAt
}

type ParticipantAddedEvent struct {
	Event
	UserID string
}

func NewEventFactory(ctx context.Context) *EventFactory {
	return &EventFactory{}
}

type EventFactory struct {
}
