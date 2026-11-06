package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeCreateRule = cqrs.CommandType("cmd.qualification.create_rule")

type CreateRuleCommand struct {
	EquipmentID     equipment.ID
	SpecificationID specifications.ID
}

func (c CreateRuleCommand) Type() cqrs.CommandType {
	return CommandTypeCreateRule
}

type CreateRuleRepositoryPort interface {
	qualification.SpecificationReader
	qualification.Reader
	qualification.Saver
}

func CreateRule(repository CreateRuleRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateRuleCommand)
		specification, err := repository.Specification(ctx, cmd.SpecificationID)
		if err != nil {
			return nil, err
		}
		qualification, err := repository.Qualification(ctx, cmd.EquipmentID)
		if err != nil {
			return nil, err
		}
		if err := qualification.CreateRule(specification); err != nil {
			return nil, err
		}
		defer qualification.ClearEvents()
		if err := repository.Save(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events(), nil
	}
}
