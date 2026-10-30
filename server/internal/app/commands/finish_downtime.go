package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeFinishDowntime = cqrs.CommandType("finish_downtime")

type FinishDowntimeCommand struct {
	DowntimeID         string    `json:"downtime_id"`
	DowntimeFinishTime time.Time `json:"downtime_finish_time"`
}

func (c FinishDowntimeCommand) Validate() error {
	return domain.DowntimeID(c.DowntimeID).Validate()
}

func (c FinishDowntimeCommand) Type() cqrs.CommandType {
	return CommandTypeFinishDowntime
}

type FinishDowntimeRepository interface {
	domain.DowntimeProvider
	domain.DowntimeSaver
}

func FinishDowntime(repository FinishDowntimeRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishDowntimeCommand)
		downtime, err := repository.Downtime(ctx,
			domain.DowntimeID(cmd.DowntimeID),
		)
		if err != nil {
			return nil, err
		}
		if err := downtime.Finish(cmd.DowntimeFinishTime); err != nil {
			return nil, err
		}
		defer downtime.Clear()
		if err := repository.SaveDowntime(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events(), nil
	}
}
