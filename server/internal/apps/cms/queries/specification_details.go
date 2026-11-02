package queries

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

var ErrSpecificationDetailsNotFound = errors.New("queries: specification details not found")

const QueryTypeSpecificationDetails = cqrs.QueryType("specification_details")

type SpecificationDetailsQuery struct {
	schema.SpecificationDetailsParams
}

func (q SpecificationDetailsQuery) Type() cqrs.QueryType {
	return QueryTypeSpecificationDetails
}

type SpecificationDetailsStore interface {
	SpecificationDetails(ctx context.Context, params schema.SpecificationDetailsParams) (schema.SpecificationDetails, error)
}

func SpecificationDetails(store SpecificationDetailsStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.SpecificationDetails(ctx, q.(SpecificationDetailsQuery).SpecificationDetailsParams)
	}
}
