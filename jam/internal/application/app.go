package application

import "jam/internal/domain/jam"

type EventPublisher interface {
	Publish(event jam.Event) error
}
