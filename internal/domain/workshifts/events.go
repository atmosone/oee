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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
}

func (e StartedEvent) Type() cqrs.EventType {
	return EventTypeStarted
}

type FinishedEvent struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   time.Duration `json:"duration"`
}

func (e FinishedEvent) Type() cqrs.EventType {
	return EventTypeFinished
}
