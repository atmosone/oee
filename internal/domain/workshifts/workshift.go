package workshifts

import (
	"errors"
	"fmt"
	"time"
)

var ErrAlreadyFinished = errors.New("workshifts: already finished")
var ErrActiveNotExists = errors.New("workshifts: active not exists")
var ErrValidationFailed = errors.New("workshifts: validation failed")

type ID string

type Name string

type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

type Workshift struct {
	ID         ID
	Name       Name
	Status     Status
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   time.Duration
}

func start(id ID, name Name, startedAt time.Time) (*Workshift, error) {
	if id == "" {
		return nil, fmt.Errorf("%w (empty id value)", ErrValidationFailed)
	}
	if name == "" {
		return nil, fmt.Errorf("%w (empty name value)", ErrValidationFailed)
	}
	if startedAt.IsZero() {
		return nil, fmt.Errorf("%w (zero start time value)", ErrValidationFailed)
	}
	return &Workshift{
		ID:        id,
		Name:      name,
		Status:    StatusActive,
		StartedAt: startedAt.UTC(),
	}, nil
}

func (w *Workshift) finish(finishedAt time.Time) error {
	if finishedAt.IsZero() {
		return fmt.Errorf("%w (zero finish time value)", ErrValidationFailed)
	}
	if w.Status == StatusFinished {
		return ErrAlreadyFinished
	}
	if w.StartedAt.After(finishedAt) {
		return fmt.Errorf("%w (invalid finish time value)", ErrValidationFailed)
	}
	w.Status = StatusFinished
	w.FinishedAt = finishedAt.UTC()
	w.Duration = w.FinishedAt.Sub(w.StartedAt)
	return nil
}
