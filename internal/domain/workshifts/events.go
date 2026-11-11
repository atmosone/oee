package workshifts

import (
	"time"

	"github.com/atmosone/cqrs"
)

const (
	EventTypeStarted  = cqrs.EventType("workshift_started")
	EventTypeFinished = cqrs.EventType("workshift_finished")
)

type StartedEvent struct {
	ID        string
	Name      string
	Status    string
	StartedAt time.Time
}

func (e StartedEvent) Type() cqrs.EventType {
	return EventTypeStarted
}

type FinishedEvent struct {
	ID         string
	Status     string
	FinishedAt time.Time
	Duration   time.Duration
}

func (e FinishedEvent) Type() cqrs.EventType {
	return EventTypeFinished
}
