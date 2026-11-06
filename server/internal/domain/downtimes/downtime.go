package downtimes

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

var ErrActiveAlreadyExists = errors.New("downtimes: active downtime already exists")
var ErrInvalidStartTime = errors.New("downtimes: invalid start time")
var ErrInvalidFinishTime = errors.New("downtimes: invalid finish time")
var ErrAlreadyFinished = errors.New("downtimes: already finished")

type ID struct{ value string }

func NewID(value string) (ID, error) {
	return ID{value: value}, nil
}

func (id ID) Value() string { return id.value }

type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

type Downtime struct {
	id          ID
	equipmentID equipment.ID
	status      Status
	startedAt   time.Time
	finishedAt  time.Time
	duration    time.Duration
	events      []cqrs.Event
}

func (d *Downtime) Events() []cqrs.Event {
	events := make([]cqrs.Event, len(d.events))
	copy(events, d.events)
	return events
}

func (d *Downtime) ClearEvents() { d.events = nil }

func (d *Downtime) Finish(finishedAt time.Time) error {
	if d.status == StatusFinished {
		return ErrAlreadyFinished
	}
	if finishedAt.IsZero() || finishedAt.Before(d.startedAt) {
		return ErrInvalidFinishTime
	}
	d.status = StatusFinished
	d.finishedAt = finishedAt.UTC()
	d.duration = d.finishedAt.Sub(d.startedAt)
	d.events = append(d.events, FinishedEvent{
		ID:         d.id.value,
		Status:     string(d.status),
		Duration:   d.duration,
		FinishedAt: finishedAt,
	})
	return nil
}

type ActiveChecker interface {
	ActiveDowntime(context.Context, equipment.ID) (bool, error)
}

type Saver interface {
	Save(context.Context, *Downtime) error
}

type Reader interface {
	Downtime(context.Context, ID) (*Downtime, error)
}

func Start(
	ctx context.Context,
	activeChecker ActiveChecker,
	id ID,
	equipmentID equipment.ID,
	startedAt time.Time,
) (*Downtime, error) {
	if startedAt.IsZero() {
		return nil, ErrInvalidStartTime
	}
	exists, err := activeChecker.ActiveDowntime(ctx, equipmentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrActiveAlreadyExists
	}
	return &Downtime{
		id:          id,
		equipmentID: equipmentID,
		status:      StatusActive,
		startedAt:   startedAt.UTC(),
		events: []cqrs.Event{StartedEvent{
			ID:          id.value,
			EquipmentID: equipmentID.Value(),
			Status:      string(StatusActive),
			StartedAt:   startedAt.UTC(),
		}},
	}, nil
}
