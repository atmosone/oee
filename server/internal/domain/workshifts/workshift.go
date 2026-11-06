package workshifts

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
)

var (
	ErrActiveAlreadyExists = errors.New("workshifts: active already exists")
	ErrAlreadyFinished     = errors.New("workshifts: already finished")
	ErrInvalidStartTime    = errors.New("workshifts: invalid start time")
	ErrInvalidFinishTime   = errors.New("workshifts: invalid finish time")
)

type ID struct{ value string }

func NewID(value string) (ID, error) {
	return ID{value: value}, nil
}

type Name struct{ value string }

func NewName(value string) (Name, error) {
	return Name{value: value}, nil
}

type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

type Workshift struct {
	id         ID
	name       Name
	status     Status
	startedAt  time.Time
	finishedAt time.Time
	duration   time.Duration
	events     []cqrs.Event
}

func (w *Workshift) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(w.events))
	copy(events, w.events)
	return events
}

func (w *Workshift) ClearEvents() { w.events = nil }

func (w *Workshift) Finish(finishedAt time.Time) error {
	if w.status == StatusFinished {
		return ErrAlreadyFinished
	}
	if finishedAt.IsZero() || finishedAt.Before(w.startedAt) {
		return ErrInvalidFinishTime
	}
	w.status = StatusFinished
	w.finishedAt = finishedAt.UTC()
	w.duration = w.finishedAt.Sub(w.startedAt)
	w.events = append(w.events, FinishedEvent{
		ID:         w.id.value,
		Status:     string(w.status),
		FinishedAt: w.finishedAt,
		Duration:   w.duration,
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
		id:        id,
		name:      name,
		status:    StatusActive,
		startedAt: startedAt.UTC(),
		events: []cqrs.Event{StartedEvent{
			ID:        id.value,
			Name:      name.value,
			Status:    string(StatusActive),
			StartedAt: startedAt.UTC(),
		}},
	}, nil
}
