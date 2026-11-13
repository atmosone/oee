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
	workshifts.Reader
	workshifts.Saver
}

func StartWorkshift(repository StartWorkshiftRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartWorkshiftCommand)
		changeover, err := repository.Changeover(ctx)
		if err != nil {
			return nil, err
		}
		if err := changeover.Start(
			cmd.ID,
			cmd.Name,
			cmd.StartedAt,
		); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, changeover); err != nil {
			return nil, err
		}
		return changeover.Events, nil
	}
}
