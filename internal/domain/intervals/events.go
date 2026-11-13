package intervals

import (
	"time"

	"github.com/atmosone/cqrs"
)

const (
	EventTypeStarted  = cqrs.EventType("interval_started")
	EventTypeFinished = cqrs.EventType("interval_finished")
)

type StartedEvent struct {
	ID          string    `json:"id"`
	EquipmentID string    `json:"equipment_id"`
	Status      string    `json:"status"`
	Category    string    `json:"category"`
	StartedAt   time.Time `json:"started_at"`
}

func (e StartedEvent) Type() cqrs.EventType {
	return EventTypeStarted
}

type FinishedEvent struct {
	ID          string        `json:"id"`
	EquipmentID string        `json:"equipment_id"`
	Status      string        `json:"status"`
	Category    string        `json:"category"`
	FinishedAt  time.Time     `json:"finished_at"`
	Duration    time.Duration `json:"duration"`
}

func (e FinishedEvent) Type() cqrs.EventType {
	return EventTypeFinished
}
