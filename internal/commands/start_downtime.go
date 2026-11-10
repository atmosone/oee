package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
)

const CommandTypeStartDowntime = cqrs.CommandType("cmd.start_downtime")

type StartDowntimeCommand struct {
	ID          downtimes.ID
	EquipmentID equipment.ID
	StartedAt   time.Time
}

func (c StartDowntimeCommand) Type() cqrs.CommandType {
	return CommandTypeStartDowntime
}

type StartDowntimeRepositoryPort interface {
	downtimes.ActiveChecker
	downtimes.Saver
}

func StartDowntime(repository StartDowntimeRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartDowntimeCommand)
		downtime, err := downtimes.Start(ctx,
			repository,
			cmd.ID,
			cmd.EquipmentID,
			cmd.StartedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events, nil
	}
}
