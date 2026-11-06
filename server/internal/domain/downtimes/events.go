package downtimes

import (
	"time"

	"github.com/atmosone/cqrs"
)

const (
	EventTypeStarted  = cqrs.EventType("evt.downtimes.started")
	EventTypeFinished = cqrs.EventType("evt.downtimes.finished")
)

type StartedEvent struct {
	ID          string
	EquipmentID string
	Status      string
	StartedAt   time.Time
}

func (e StartedEvent) Type() cqrs.EventType {
	return EventTypeStarted
}

type FinishedEvent struct {
	ID         string
	Status     string
	Duration   time.Duration
	FinishedAt time.Time
}

func (e FinishedEvent) Type() cqrs.EventType {
	return EventTypeFinished
}
