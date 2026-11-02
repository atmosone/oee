package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeSpecificationsList = cqrs.QueryType("specifications_list")

type SpecificationsListQuery struct{}

func (q SpecificationsListQuery) Type() cqrs.QueryType { return QueryTypeSpecificationsList }

type SpecificationsListStore interface {
	SpecificationsList(ctx context.Context) (schema.SpecificationsList, error)
}

func SpecificationsList(store SpecificationsListStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.SpecificationsList(ctx)
	}
}
