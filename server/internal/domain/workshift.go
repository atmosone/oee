package domain

import (
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
)

type WorkshiftID string

func (wid WorkshiftID) Validate() error {
	if wid == "" {
		return fmt.Errorf(
			"%w (empty workshift id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type WorkshiftName string

func (wn WorkshiftName) Validate() error {
	const min, max = 1, 64
	if l := len(wn); l < min || l > max {
		return fmt.Errorf(
			"%w (workshift name can contains from %d to %d characters)",
			ErrValidationFailed, min, max)
	}
	return nil
}

type WorkshiftStatus string

const (
	WorkshiftStatusActive   = WorkshiftStatus("active")
	WorkshiftStatusFinished = WorkshiftStatus("finished")
)

type Workshift struct {
	*reactiveEntity
	id         WorkshiftID
	name       WorkshiftName
	status     WorkshiftStatus
	duration   time.Duration
	startedAt  time.Time
	finishedAt time.Time
}

const EventTypeWorkshiftStarted = cqrs.EventType("workshift_started")

type WorkshiftStartedEvent struct {
	WorkshiftID        string    `json:"workshift_id"`
	WorkshiftName      string    `json:"workshift_name"`
	WorkshiftStatus    string    `json:"workshift_status"`
	WorkshiftStartedAt time.Time `json:"workshift_started_at"`
}

func (e WorkshiftStartedEvent) Type() cqrs.EventType {
	return EventTypeWorkshiftStarted
}

func NewWorkshift(
	id WorkshiftID,
	name WorkshiftName,
	startedAt time.Time,
) *Workshift {
	w := &Workshift{
		reactiveEntity: &reactiveEntity{},
		id:             id,
		name:           name,
		status:         WorkshiftStatusActive,
		startedAt:      startedAt.UTC(),
	}
	w.events = append(w.events, WorkshiftStartedEvent{
		WorkshiftID:        string(w.id),
		WorkshiftName:      string(w.name),
		WorkshiftStatus:    string(w.status),
		WorkshiftStartedAt: w.startedAt,
	})
	return w
}

const EventTypeWorkshiftFinished = cqrs.EventType("workshift_finished")

type WorkshiftFinishedEvent struct {
	WorkshiftID         string        `json:"workshift_id"`
	WorkshiftStatus     string        `json:"workshift_status"`
	WorkshiftDuration   time.Duration `json:"workshift_duration"`
	WorkshiftFinishedAt time.Time     `json:"workshift_finished_at"`
}

func (e WorkshiftFinishedEvent) Type() cqrs.EventType {
	return EventTypeWorkshiftFinished
}

func (w *Workshift) Finish(finishedAt time.Time) error {
	if err := ValidateTimestamp(finishedAt); err != nil {
		return err
	}
	if finishedAt.Before(w.startedAt) {
		return fmt.Errorf(
			"%w (workshift finish time should be after start time)",
			ErrUnprocessable,
		)
	}
	if w.status == WorkshiftStatusFinished {
		return fmt.Errorf(
			"%w (workshift already finished)",
			ErrConflict,
		)
	}
	w.status = WorkshiftStatusFinished
	w.finishedAt = finishedAt.UTC()
	w.duration = w.finishedAt.Sub(w.startedAt)
	w.events = append(w.events, WorkshiftFinishedEvent{
		WorkshiftID:         string(w.id),
		WorkshiftStatus:     string(w.status),
		WorkshiftDuration:   w.duration,
		WorkshiftFinishedAt: w.finishedAt,
	})
	return nil
}
