package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
)

const CommandTypeStart = cqrs.CommandType("cmd.downtimes.start")

type StartCommand struct {
	ID          downtimes.ID
	EquipmentID equipment.ID
	StartedAt   time.Time
}

func (c StartCommand) Type() cqrs.CommandType {
	return CommandTypeStart
}

type StartRepositoryPort interface {
	downtimes.ActiveChecker
	downtimes.Saver
}

func Start(repository StartRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartCommand)
		downtime, err := downtimes.Start(ctx,
			repository,
			cmd.ID,
			cmd.EquipmentID,
			cmd.StartedAt,
		)
		if err != nil {
			return nil, err
		}
		defer downtime.ClearEvents()
		if err := repository.Save(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events(), nil
	}
}
