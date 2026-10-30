package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeExpandEquipmentLine = cqrs.CommandType("expand_equipment_line")

type ExpandEquipmentLineCommand struct {
	EquipmentLineID string `json:"equipment_line_id"`
	EquipmentID     string `json:"equipment_id"`
}

func (c ExpandEquipmentLineCommand) Type() cqrs.CommandType {
	return CommandTypeExpandEquipmentLine
}

func (c ExpandEquipmentLineCommand) Validate() error {
	return errors.Join(
		domain.EquipmentLineID(c.EquipmentLineID).Validate(),
		domain.EquipmentID(c.EquipmentID).Validate(),
	)
}

type ExpandEquipmentLineRepository interface {
	domain.EquipmentLineProvider
	domain.EquipmentLineSaver
}

func ExpandEquipmentLine(repository ExpandEquipmentLineRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(ExpandEquipmentLineCommand)
		equipmentLine, err := repository.EquipmentLine(ctx,
			domain.EquipmentLineID(cmd.EquipmentLineID),
		)
		if err != nil {
			return nil, err
		}
		if err := equipmentLine.Expand(
			domain.EquipmentID(cmd.EquipmentID),
		); err != nil {
			return nil, err
		}
		defer equipmentLine.Clear()
		if err := repository.SaveEquipmentLine(ctx, equipmentLine); err != nil {
			return nil, err
		}
		return equipmentLine.Events(), nil
	}
}
