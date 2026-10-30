package commands

import (
	"context"
	"errors"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeStartDowntime = cqrs.CommandType("start_downtime")

type StartDowntimeCommand struct {
	DowntimeID        string    `json:"downtime_id"`
	EquipmentID       string    `json:"equipment_id"`
	DowntimeStartTime time.Time `json:"downtime_start_time"`
}

func (c StartDowntimeCommand) Validate() error {
	return errors.Join(
		domain.DowntimeID(c.DowntimeID).Validate(),
		domain.EquipmentID(c.EquipmentID).Validate(),
	)
}

func (c StartDowntimeCommand) Type() cqrs.CommandType {
	return CommandTypeStartDowntime
}

type StartDowntimeRepository interface {
	domain.DowntimeSaver
}

func StartDowntime(repository StartDowntimeRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartDowntimeCommand)
		downtime := domain.NewDowntime(
			domain.DowntimeID(cmd.DowntimeID),
			domain.EquipmentID(cmd.EquipmentID),
			cmd.DowntimeStartTime,
		)
		defer downtime.Clear()
		if err := repository.SaveDowntime(ctx, downtime); err != nil {
			return nil, err
		}
		return downtime.Events(), nil
	}
}
