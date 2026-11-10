package workshifts

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
)

var (
	ErrWorkshiftNotFound   = errors.New("workshifts: workshift not found")
	ErrActiveAlreadyExists = errors.New("workshifts: active already exists")
	ErrAlreadyFinished     = errors.New("workshifts: already finished")
	ErrInvalidStartTime    = errors.New("workshifts: invalid start time")
	ErrInvalidFinishTime   = errors.New("workshifts: invalid finish time")
)

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
	Events     []cqrs.Event
}

func (w *Workshift) Finish(finishedAt time.Time) error {
	if w.Status == StatusFinished {
		return ErrAlreadyFinished
	}
	if finishedAt.IsZero() || finishedAt.Before(w.StartedAt) {
		return ErrInvalidFinishTime
	}
	w.Status = StatusFinished
	w.FinishedAt = finishedAt.UTC()
	w.Duration = w.FinishedAt.Sub(w.StartedAt)
	w.Events = append(w.Events, FinishedEvent{
		ID:         string(w.ID),
		Status:     string(w.Status),
		FinishedAt: w.FinishedAt,
		Duration:   w.Duration,
	})
	return nil
}

type ActiveChecker interface {
	ActiveWorkshift(context.Context) (bool, error)
}

type Saver interface {
	Save(context.Context, *Workshift) error
}

type Reader interface {
	Workshift(context.Context, ID) (*Workshift, error)
}

func Start(
	ctx context.Context,
	activeChecker ActiveChecker,
	id ID,
	name Name,
	startedAt time.Time,
) (*Workshift, error) {
	if startedAt.IsZero() {
		return nil, ErrInvalidStartTime
	}
	exists, err := activeChecker.ActiveWorkshift(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrActiveAlreadyExists
	}
	return &Workshift{
		ID:        id,
		Name:      name,
		Status:    StatusActive,
		StartedAt: startedAt.UTC(),
		Events: []cqrs.Event{StartedEvent{
			ID:        string(id),
			Name:      string(name),
			Status:    string(StatusActive),
			StartedAt: startedAt.UTC(),
		}},
	}, nil
}
