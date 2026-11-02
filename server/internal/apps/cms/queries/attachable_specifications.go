package queries

import (
	"context"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

const QueryTypeAttachableSpecifications = cqrs.QueryType("attachable_specifications")

type AttachableSpecificationsQuery struct {
	schema.AttachableSpecificationsParams
}

func (q AttachableSpecificationsQuery) Type() cqrs.QueryType {
	return QueryTypeAttachableSpecifications
}

type AttachableSpecificationsStore interface {
	AttachableSpecifications(ctx context.Context, params schema.AttachableSpecificationsParams) (schema.AttachableSpecificationsList, error)
}

func AttachableSpecifications(store AttachableSpecificationsStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.AttachableSpecifications(ctx, q.(AttachableSpecificationsQuery).AttachableSpecificationsParams)
	}
}
