package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeCreateQualificationRule = cqrs.CommandType("create_qualification_rule")

type CreateQualificationRuleCommand struct {
	EquipmentID     string `json:"equipment_id"`
	SpecificationID string `json:"specification_id"`
}

func (c CreateQualificationRuleCommand) Validate() error {
	return errors.Join(
		domain.EquipmentID(c.EquipmentID).Validate(),
		domain.SpecificationID(c.SpecificationID).Validate(),
	)
}

func (c CreateQualificationRuleCommand) Type() cqrs.CommandType {
	return CommandTypeCreateQualificationRule
}

type CreateQualificationRuleRepository interface {
	domain.QualificationRuleSaver
}

func CreateQualificationRule(repository CreateQualificationRuleRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateQualificationRuleCommand)
		qualificationRule := domain.NewQualificationRule(
			domain.EquipmentID(cmd.EquipmentID),
			domain.SpecificationID(cmd.SpecificationID),
		)
		defer qualificationRule.Clear()
		if err := repository.SaveQualificationRule(ctx, qualificationRule); err != nil {
			return nil, err
		}
		return qualificationRule.Events(), nil
	}
}
