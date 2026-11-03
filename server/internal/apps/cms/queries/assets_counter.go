package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeAssetsCounter = cqrs.QueryType("assets_counter")

type AssetsCounterQuery struct{}

func (q AssetsCounterQuery) Type() cqrs.QueryType { return QueryTypeAssetsCounter }

type AssetsCounterStore interface {
	AssetsCounter(ctx context.Context) (schema.AssetsCounter, error)
}

func AssetsCounter(store AssetsCounterStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.AssetsCounter(ctx)
	}
}
