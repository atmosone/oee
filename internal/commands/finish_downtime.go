package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
)

const CommandTypeFinishDowntime = cqrs.CommandType("finish_downtime")

type FinishDowntimeCommand struct {
	ID         downtimes.ID `json:"id"`
	FinishedAt time.Time    `json:"finished_at"`
}

func (c FinishDowntimeCommand) Type() cqrs.CommandType {
	return CommandTypeFinishDowntime
}

type FinishDowntimeRepositoryPort interface {
	downtimes.Reader
	downtimes.Saver
}

func FinishDowntime(repository FinishDowntimeRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishDowntimeCommand)
		downtime, err := repository.Downtime(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := downtime.Finish(cmd.FinishedAt); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events, nil
	}
}
