package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"

	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/intervals"
)

const CommandTypeStartInterval = cqrs.CommandType("start_interval")

type StartIntervalCommand struct {
	ID          intervals.ID       `json:"id"`
	EquipmentID equipment.ID       `json:"equipment_id"`
	Category    intervals.Category `json:"category"`
	StartedAt   time.Time          `json:"started_at"`
}

func (c StartIntervalCommand) Type() cqrs.CommandType {
	return CommandTypeStartInterval
}

type StartIntervalRepositoryPort interface {
	intervals.Reader
	intervals.Saver
}

func StartInterval(repository StartIntervalRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(StartIntervalCommand)
		equipment, err := repository.Equipment(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := equipment.StartInterval(
			cmd.ID,
			cmd.Category,
			cmd.StartedAt,
		); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, equipment); err != nil {
			return nil, err
		}
		return equipment.Events, nil
	}
}
