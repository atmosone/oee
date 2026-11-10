package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeMoveEquipmentInLine = cqrs.CommandType("cmd.lines.move_equipment")

type MoveEquipmentInLineCommand struct {
	ID          lines.ID
	EquipmentID equipment.ID
	Direction   lines.Direction
}

func (c MoveEquipmentInLineCommand) Type() cqrs.CommandType {
	return CommandTypeMoveEquipmentInLine
}

type MoveEquipmentInLineRepositoryPort interface {
	lines.Reader
	lines.Saver
}

func MoveEquipmentInLine(repository MoveEquipmentInLineRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(MoveEquipmentInLineCommand)
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
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events, nil
	}
}
