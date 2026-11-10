package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

const CommandTypeFinishWorkshift = cqrs.CommandType("cmd.workshifts.finish")

type FinishWorkshiftCommand struct {
	ID                  workshifts.ID
	FinishWorkshiftedAt time.Time
}

func (c FinishWorkshiftCommand) Type() cqrs.CommandType {
	return CommandTypeFinishWorkshift
}

type FinishWorkshiftRepositoryPort interface {
	workshifts.Reader
	workshifts.Saver
}

func FinishWorkshift(repository FinishWorkshiftRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishWorkshiftCommand)
		workshift, err := repository.Workshift(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := workshift.Finish(cmd.FinishWorkshiftedAt); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events, nil
	}
}
