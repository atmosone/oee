package workshifts

import (
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
)

var ErrValidation = errors.New("workshifts: validation failed")

type Workshift struct {
	id         ID
	name       Name
	status     Status
	duration   time.Duration
	startedAt  time.Time
	finishedAt time.Time
	events     []cqrs.Event
}

func Start(id ID, name Name, start time.Time) (*Workshift, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	if err := name.Validate(); err != nil {
		return nil, err
	}
	if start.IsZero() {
		return nil, fmt.Errorf("%w (zero start time)", ErrValidation)
	}
	w := &Workshift{
		id:        id,
		name:      name,
		status:    StatusActive,
		startedAt: start,
	}
	w.events = append(w.events, StartedEvent{
		ID:        string(w.id),
		Name:      string(w.name),
		Status:    string(w.status),
		StartedAt: w.startedAt,
	})
	return w, nil
}

func (w *Workshift) Finish(finish time.Time) error {
	if finish.IsZero() {
		return fmt.Errorf("%w (zero finish time)", ErrValidation)
	}
	if finish.Before(w.startedAt) {
		return fmt.Errorf("%w (unprocessable finish time)", ErrValidation)
	}
	w.finishedAt = finish
	w.duration = w.finishedAt.Sub(w.startedAt)
	w.status = StatusFinished
	w.events = append(w.events, FinishedEvent{
		ID:         string(w.id),
		Status:     string(w.status),
		Duration:   w.duration,
		FinishedAt: w.finishedAt,
	})
	return nil
}
