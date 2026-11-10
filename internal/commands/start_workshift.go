package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

const CommandTypeStartWorkshift = cqrs.CommandType("start_workshift")

type StartWorkshiftCommand struct {
	ID        workshifts.ID   `json:"id"`
	Name      workshifts.Name `json:"name"`
	StartedAt time.Time       `json:"started_at"`
}

func (c StartWorkshiftCommand) Type() cqrs.CommandType {
	return CommandTypeStartWorkshift
}

type StartWorkshiftRepositoryPort interface {
	workshifts.ActiveChecker
	workshifts.Saver
}

func StartWorkshift(repository StartWorkshiftRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartWorkshiftCommand)
		workshift, err := workshifts.Start(ctx,
			repository,
			cmd.ID,
			cmd.Name,
			cmd.StartedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, workshift); err != nil {
			return nil, err
		}
		return workshift.Events, nil
	}
}
