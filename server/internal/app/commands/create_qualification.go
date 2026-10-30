package commands

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain"
)

const CommandTypeCreateQualification = cqrs.CommandType("create_qualification")

type CreateQualificationCommand struct {
	DowntimeID           string `json:"downtime_id"`
	SpecificationID      string `json:"specification_id"`
	QualificationComment string `json:"qualification_comment"`
}

func (c CreateQualificationCommand) Validate() error {
	return errors.Join(
		domain.DowntimeID(c.DowntimeID).Validate(),
		domain.SpecificationID(c.SpecificationID).Validate(),
		domain.QualificationComment(c.QualificationComment).Validate(),
	)
}

func (c CreateQualificationCommand) Type() cqrs.CommandType {
	return CommandTypeCreateQualification
}

type CreateQualificationRepository interface {
	domain.QualificationSaver
}

func CreateQualification(repository CreateQualificationRepository) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateQualificationCommand)
		qualification := domain.NewQualification(
			domain.DowntimeID(cmd.DowntimeID),
			domain.SpecificationID(cmd.SpecificationID),
			domain.QualificationComment(cmd.QualificationComment),
		)
		defer qualification.Clear()
		if err := repository.SaveQualification(ctx, qualification); err != nil {
			return nil, err
		}
		return qualification.Events(), nil
	}
}
