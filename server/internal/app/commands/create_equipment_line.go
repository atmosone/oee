package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeCreateEquipmentLine = cqrs.CommandType("create_equipment_line")

type CreateEquipmentLineCommand struct {
	EquipmentLineID   string `json:"equipment_line_id"`
	EquipmentLineName string `json:"equipment_line_name"`
}

func (c CreateEquipmentLineCommand) Type() cqrs.CommandType {
	return CommandTypeCreateEquipmentLine
}

func (c CreateEquipmentLineCommand) Validate() error {
	return errors.Join(
		domain.EquipmentLineID(c.EquipmentLineID).Validate(),
		domain.EquipmentLineName(c.EquipmentLineName).Validate(),
	)
}

type CreateEquipmentLineRepository interface {
	domain.EquipmentLineSaver
}

func CreateEquipmentLine(repository CreateEquipmentLineRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateEquipmentLineCommand)
		equipmentLine := domain.NewEquipmentLine(
			domain.EquipmentLineID(cmd.EquipmentLineID),
			domain.EquipmentLineName(cmd.EquipmentLineName),
		)
		defer equipmentLine.Clear()
		if err := repository.SaveEquipmentLine(ctx, equipmentLine); err != nil {
			return nil, err
		}
		return equipmentLine.Events(), nil
	}
}
