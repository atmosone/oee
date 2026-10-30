package commands

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeFinishWorkshift = cqrs.CommandType("finish_workshift")

type FinishWorkshiftCommand struct {
	WorkshiftID         string    `json:"workshift_id"`
	WorkshiftFinishTime time.Time `json:"workshift_finish_time"`
}

func (c FinishWorkshiftCommand) Validate() error {
	return errors.Join(
		domain.WorkshiftID(c.WorkshiftID).Validate(),
	)
}

func (c FinishWorkshiftCommand) Type() cqrs.CommandType {
	return CommandTypeFinishWorkshift
}

type FinishWorkshiftRepository interface {
	domain.WorkshiftProvider
	domain.WorkshiftSaver
}

func FinishWorkshift(repository FinishWorkshiftRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishWorkshiftCommand)
		workshift, err := repository.Workshift(ctx,
			domain.WorkshiftID(cmd.WorkshiftID),
		)
		if err != nil {
			return nil, err
		}
		defer workshift.Clear()
		if err := workshift.Finish(cmd.WorkshiftFinishTime); err != nil {
			return nil, err
		}
		if err := repository.SaveWorkshift(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events(), nil
	}
}
