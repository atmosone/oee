package workshifts

import (
	"time"

	"github.com/atmosone/cqrs"
)

const EventTypeStarted = cqrs.EventType("workshifts.started")

type StartedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
}

func (e StartedEvent) Type() cqrs.EventType { return EventTypeStarted }

const EventTypeFinished = cqrs.EventType("workshifts.finished")

type FinishedEvent struct {
	ID         string        `json:"id"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	FinishedAt time.Time     `json:"finished_at"`
}

func (e FinishedEvent) Type() cqrs.EventType { return EventTypeFinished }
