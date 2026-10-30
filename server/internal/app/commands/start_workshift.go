package commands

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeStartWorkshift = cqrs.CommandType("start_workshift")

type StartWorkshiftCommand struct {
	WorkshiftID        string    `json:"workshift_id"`
	WorkshiftName      string    `json:"workshift_name"`
	WorkshiftStartTime time.Time `json:"workshift_start_time"`
}

func (c StartWorkshiftCommand) Validate() error {
	return errors.Join(
		domain.WorkshiftID(c.WorkshiftID).Validate(),
		domain.WorkshiftName(c.WorkshiftName).Validate(),
	)
}

func (c StartWorkshiftCommand) Type() cqrs.CommandType {
	return CommandTypeStartWorkshift
}

type StartWorkshiftRepository interface {
	domain.WorkshiftSaver
}

func StartWorkshift(repository StartWorkshiftRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartWorkshiftCommand)
		workshift := domain.NewWorkshift(
			domain.WorkshiftID(cmd.WorkshiftID),
			domain.WorkshiftName(cmd.WorkshiftName),
			cmd.WorkshiftStartTime,
		)
		defer workshift.Clear()
		if err := repository.SaveWorkshift(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events(), nil
	}
}
