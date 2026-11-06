package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeRemoveRule = cqrs.CommandType("cmd.qualification.remove_rule")

type RemoveRuleCommand struct {
	EquipmentID     equipment.ID
	SpecificationID specifications.ID
}

func (c RemoveRuleCommand) Type() cqrs.CommandType {
	return CommandTypeRemoveRule
}

type RemoveRuleRepositoryPort interface {
	qualification.Reader
	qualification.Saver
}

func RemoveRule(repository RemoveRuleRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(RemoveRuleCommand)
		qualification, err := repository.Qualification(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := qualification.RemoveRule(cmd.SpecificationID); err != nil {
			return nil, err
		}
		defer qualification.ClearEvents()
		if err := repository.Save(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events(), nil
	}
}
