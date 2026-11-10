package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeDeprecateSpecification = cqrs.CommandType("deprecate_specification")

type DeprecateSpecificationCommand struct {
	ID specifications.ID `json:"id"`
}

func (c DeprecateSpecificationCommand) Type() cqrs.CommandType {
	return CommandTypeDeprecateSpecification
}

type DeprecateSpecificationRepositoryPort interface {
	specifications.Reader
	specifications.Saver
}

func DeprecateSpecification(repository DeprecateSpecificationRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(DeprecateSpecificationCommand)
		specification, err := repository.Specification(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := specification.Deprecate(); err != nil {
			return nil, err
		}
		if err := repository.Save(ctx, specification); err != nil {
			return nil, err
		}
		return specification.Events, nil
	}
}
