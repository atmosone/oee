package downtimes

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
)

var ErrDowntimeNotFound = errors.New("downtimes: downtime not found")
var ErrActiveAlreadyExists = errors.New("downtimes: active downtime already exists")
var ErrInvalidStartTime = errors.New("downtimes: invalid start time")
var ErrInvalidFinishTime = errors.New("downtimes: invalid finish time")
var ErrAlreadyFinished = errors.New("downtimes: already finished")

// ID - downtime ID.
type ID string

// Status - downtime status.
type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

type Downtime struct {
	ID          ID
	EquipmentID equipment.ID
	Status      Status
	StartedAt   time.Time
	FinishedAt  time.Time
	Duration    time.Duration
	Events      []cqrs.Event
}

func (d *Downtime) Finish(finishedAt time.Time) error {
	if d.Status == StatusFinished {
		return ErrAlreadyFinished
	}
	if finishedAt.IsZero() || finishedAt.Before(d.StartedAt) {
		return ErrInvalidFinishTime
	}
	d.Status = StatusFinished
	d.FinishedAt = finishedAt.UTC()
	d.Duration = d.FinishedAt.Sub(d.StartedAt)
	d.Events = append(d.Events, FinishedEvent{
		ID:         string(d.ID),
		Status:     string(d.Status),
		Duration:   d.Duration,
		FinishedAt: d.FinishedAt,
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
		ID:          ID(id),
		EquipmentID: equipmentID,
		Status:      StatusActive,
		StartedAt:   startedAt.UTC(),
		Events: []cqrs.Event{StartedEvent{
			ID:          string(id),
			EquipmentID: string(equipmentID),
			Status:      string(StatusActive),
			StartedAt:   startedAt.UTC(),
		}},
	}, nil
}
