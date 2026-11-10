package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeMoveEquipmentInLine = cqrs.CommandType("move_equipment_in_line")

type MoveEquipmentInLineCommand struct {
	ID          lines.ID        `json:"id"`
	EquipmentID equipment.ID    `json:"equipment_id"`
	Direction   lines.Direction `json:"direction"`
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
