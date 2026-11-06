package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

const CommandTypeStart = cqrs.CommandType("cmd.workshifts.start")

type StartCommand struct {
	ID        workshifts.ID
	Name      workshifts.Name
	StartedAt time.Time
}

func (c StartCommand) Type() cqrs.CommandType {
	return CommandTypeStart
}

type StartRepositoryPort interface {
	workshifts.ActiveChecker
	workshifts.Saver
}

func Start(repository StartRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartCommand)
		workshift, err := workshifts.Start(ctx,
			repository,
			cmd.ID,
			cmd.Name,
			cmd.StartedAt,
		)
		if err != nil {
			return nil, err
		}
		defer workshift.ClearEvents()
		if err := repository.Save(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events(), nil
	}
}
