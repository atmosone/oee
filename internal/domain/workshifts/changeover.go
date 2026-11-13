package workshifts

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
)

// Changeover - workshift management aggregate root.
type Changeover struct {
	Active *Workshift
	Events []cqrs.Event
}

// Start - finish previous workshift by next workshift start time and start new workshift.
func (c *Changeover) Start(id ID, name Name, startedAt time.Time) error {
	if err := c.Finish(startedAt); err != nil &&
		!errors.Is(err, ErrActiveNotExists) {
		return err
	}
	active, err := start(id, name, startedAt)
	if err != nil {
		return err
	}
	c.Active = active
	c.Events = append(c.Events, StartedEvent{
		ID:        string(c.Active.ID),
		Name:      string(c.Active.Name),
		Status:    string(c.Active.Status),
		StartedAt: c.Active.StartedAt,
	})
	return nil
}

// Finish - finish active workshift.
func (c *Changeover) Finish(finishedAt time.Time) error {
	if c.Active == nil {
		return ErrActiveNotExists
	}
	if err := c.Active.finish(finishedAt); err != nil {
		return err
	}
	c.Events = append(c.Events, FinishedEvent{
		ID:         string(c.Active.ID),
		Status:     string(c.Active.Status),
		FinishedAt: c.Active.FinishedAt,
		Duration:   c.Active.Duration,
	})
	c.Active = nil
	return nil
}

type Saver interface {
	Save(context.Context, *Changeover) error
}

type Reader interface {
	Changeover(context.Context) (*Changeover, error)
}
