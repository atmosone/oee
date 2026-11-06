package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

const CommandTypeFinish = cqrs.CommandType("cmd.workshifts.finish")

type FinishCommand struct {
	ID         workshifts.ID
	FinishedAt time.Time
}

func (c FinishCommand) Type() cqrs.CommandType {
	return CommandTypeFinish
}

type FinishRepositoryPort interface {
	workshifts.Reader
	workshifts.Saver
}

func Finish(repository FinishRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishCommand)
		workshift, err := repository.Workshift(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := workshift.Finish(cmd.FinishedAt); err != nil {
			return nil, err
		}
		defer workshift.ClearEvents()
		if err := repository.Save(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events(), nil
	}
}
