package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/specifications"
)

const CommandTypeDeprecate = cqrs.CommandType("cmd.specifications.deprecate")

type DeprecateCommand struct {
	ID specifications.ID
}

func (c DeprecateCommand) Type() cqrs.CommandType {
	return CommandTypeDeprecate
}

type DeprecateRepositoryPort interface {
	specifications.Reader
	specifications.Saver
}

func Deprecate(repository DeprecateRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(DeprecateCommand)
		specification, err := repository.Specification(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if err := specification.Deprecate(); err != nil {
			return nil, err
		}
		defer specification.ClearEvents()
		if err := repository.Save(ctx, specification); err != nil {
			return nil, err
		}
		return specification.Events(), nil
	}
}
