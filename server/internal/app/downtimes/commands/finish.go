package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
)

const CommandTypeFinish = cqrs.CommandType("cmd.downtimes.finish")

type FinishCommand struct {
	ID         downtimes.ID
	FinishedAt time.Time
}

func (c FinishCommand) Type() cqrs.CommandType {
	return CommandTypeFinish
}

type FinishRepositoryPort interface {
	downtimes.Reader
	downtimes.Saver
}

func Finish(repository FinishRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishCommand)
		downtime, err := repository.Downtime(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := downtime.Finish(cmd.FinishedAt); err != nil {
			return nil, err
		}
		defer downtime.ClearEvents()
		if err := repository.Save(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events(), nil
	}
}
