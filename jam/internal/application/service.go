package application

import (
	"context"
	"jam/internal/domain/jam"
	"time"
)

type EventPublisher interface {
	Publish(ctx context.Context, event jam.DomainEvent) error
}

func NewJamService(
	eventBus EventPublisher,
	jamRepository jam.JamRepository,
	participantRepository jam.ParticipantRepository,
) *JamService {
	return &JamService{
		eventBus:              eventBus,
		jamRepository:         jamRepository,
		participantRepository: participantRepository,
	}
}

type JamService struct {
	eventBus              EventPublisher
	jamRepository         jam.JamRepository
	participantRepository jam.ParticipantRepository
}

type CreateJamParams struct {
	Name      string
	Location  string
	Date      time.Time
	CreatedBy string
}

func (js *JamService) CreateJam(ctx context.Context, params CreateJamParams) error {
	j := jam.NewJam(params.Name, params.Location, time.Now(), jam.UserID(params.CreatedBy))
	err := js.jamRepository.Store(&j)
	if err != nil {
		return err
	}

	err = js.eventBus.Publish(ctx, &jam.JamCreatedEvent{
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

func (js *JamService) ListJams(ctx context.Context) {
	js.jamRepository.FindAll()
}
