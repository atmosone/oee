package queries

import (
	"context"
	"errors"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

var ErrEquipmentDetailsNotFound = errors.New("queries: equipment details not found")

const QueryTypeEquipmentDetails = cqrs.QueryType("equipment_details")

type EquipmentDetailsQuery struct {
	schema.EquipmentDetailsParams
}

func (q EquipmentDetailsQuery) Type() cqrs.QueryType {
	return QueryTypeEquipmentDetails
}

type EquipmentDetailsStore interface {
	EquipmentDetails(ctx context.Context, params schema.EquipmentDetailsParams) (schema.EquipmentDetails, error)
}

func EquipmentDetails(store EquipmentDetailsStore) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		return store.EquipmentDetails(ctx, q.(EquipmentDetailsQuery).EquipmentDetailsParams)
	}
}
