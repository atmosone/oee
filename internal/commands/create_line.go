package commands

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/lines"
)

const CommandTypeCreateLine = cqrs.CommandType("cmd.lines.create")

type CreateLineCommand struct {
	ID   lines.ID
	Name lines.Name
}

func (c CreateLineCommand) Type() cqrs.CommandType {
	return CommandTypeCreateLine
}

type CreateLineRepositoryPort interface {
	lines.Saver
}

func CreateLine(repository CreateLineRepositoryPort) cqrs.CommandHandlerFunc {
	return func(ctx context.Context, c cqrs.Command) ([]cqrs.Event, error) {
		cmd := c.(CreateLineCommand)
		line := lines.New(cmd.ID, cmd.Name)
		if err := repository.Save(ctx, line); err != nil {
			return nil, err
		}
		return line.Events, nil
	}
}
