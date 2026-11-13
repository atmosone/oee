package commands

import (
	"context"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/intervals"
)

const CommandTypeFinishInterval = cqrs.CommandType("finish_interval")

type FinishIntervalCommand struct {
	EquipmentID equipment.ID `json:"equipment_id"`
	FinishedAt  time.Time    `json:"finished_at"`
}

func (c FinishIntervalCommand) Type() cqrs.CommandType {
	return CommandTypeFinishInterval
}

type FinishIntervalRepositoryPort interface {
	intervals.Reader
	intervals.Saver
}

func FinishInterval(repository FinishIntervalRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(FinishIntervalCommand)
		equipment, err := repository.Equipment(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := equipment.Finish(cmd.FinishedAt); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, equipment); err != nil {
			return nil, err
		}
		return equipment.Events, nil
	}
}
