package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeMoveEquipment = cqrs.CommandType("cmd.lines.move_equipment")

type MoveEquipmentCommand struct {
	ID          lines.ID
	EquipmentID equipment.ID
	Direction   lines.Direction
}

func (c MoveEquipmentCommand) Type() cqrs.CommandType {
	return CommandTypeMoveEquipment
}

type MoveEquipmentRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func MoveEquipment(repository MoveEquipmentRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(MoveEquipmentCommand)
		line, err := repository.Line(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := line.MoveEquipment(
			cmd.EquipmentID,
			cmd.Direction,
		); err != nil {
			return nil, err
		}
		defer line.ClearEvents()
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events(), nil
	}
}
