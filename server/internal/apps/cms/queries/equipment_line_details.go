package queries

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

var ErrEquipmentLineDetailsNotFound = errors.New("queries: equipment line details not found")

const QueryTypeEquipmentLineDetails = cqrs.QueryType("equipment_line_details")

type EquipmentLineDetailsQuery struct {
	schema.EquipmentLineDetailsParams
}

func (q EquipmentLineDetailsQuery) Type() cqrs.QueryType {
	return QueryTypeEquipmentLineDetails
}

type EquipmentLineDetailsStore interface {
	EquipmentLineDetails(ctx context.Context, params schema.EquipmentLineDetailsParams) (schema.EquipmentLineDetails, error)
}

func EquipmentLineDetails(store EquipmentLineDetailsStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.EquipmentLineDetails(ctx, q.(EquipmentLineDetailsQuery).EquipmentLineDetailsParams)
	}
}
