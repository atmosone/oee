package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

const CommandTypeFinishWorkshift = cqrs.CommandType("finish_workshift")

type FinishWorkshiftCommand struct {
	FinishedAt time.Time `json:"finished_at"`
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
		changeover, err := repository.Changeover(ctx)
		if err != nil {
			return nil, err
		}
		if err := changeover.Finish(cmd.FinishedAt); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, changeover); err != nil {
			return nil, err
		}
		return changeover.Events, nil
	}
}
