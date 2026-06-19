package application

import (
	"context"
	"jam/internal/domain/jam"
	"time"
)

type EventPublisher interface {
	Publish(event jam.DomainEvent) error
}

func NewJamService(eventBus EventPublisher) *JamService {
	return &JamService{
		eventBus: eventBus,
	}
}

type JamService struct {
	eventBus EventPublisher
}

type CreateJamParams struct {
	Name      string
	Location  string
	Date      time.Time
	CreatedBy string
}

func (js *JamService) CreateJam(ctx context.Context, params CreateJamParams) error {
	err := js.eventBus.Publish(jam.JamCreatedEvent{
		Event: jam.Event{
			Type:       jam.JamCreated,
			CreatedBy:  params.CreatedBy,
			OccurredAt: time.Now(),
		},
		Name:     params.Name,
		Location: params.Location,
		Date:     params.Date,
	})
	return err
}
