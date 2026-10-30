package domain

import (
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
)

type EquipmentID string

func (eid EquipmentID) Validate() error {
	if eid == "" {
		return fmt.Errorf(
			"%w (empty equipment id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type EquipmentName string

func (en EquipmentName) Validate() error {
	const min, max = 1, 64
	if l := len(en); l < min || l > max {
		return fmt.Errorf(
			"%w (equipment name can contains from %d to %d characters)",
			ErrValidationFailed, min, max,
		)
	}
	return nil
}

type EquipmentUnits string

func (eu EquipmentUnits) Validate() error {
	const min, max = 1, 25
	if l := len(eu); l < min || l > max {
		return fmt.Errorf(
			"%w (equipment units can contains from %d to %d characters)",
			ErrValidationFailed, min, max,
		)
	}
	return nil
}

type Equipment struct {
	*reactiveEntity
	id    EquipmentID
	name  EquipmentName
	units EquipmentUnits
}

const EventTypeEquipmentCreated = cqrs.EventType("equipment_created")

type EquipmentLineCreatedEvent struct {
	EquipmentLineID        string    `json:"equipment_line_id"`
	EquipmentLineName      string    `json:"equipment_line_name"`
	EquipmentLineCreatedAt time.Time `json:"equipment_line_created_at"`
}

func (e EquipmentLineCreatedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentLineCreated
}

func NewEquipment(id EquipmentID, name EquipmentName, units EquipmentUnits) *Equipment {
	e := &Equipment{
		reactiveEntity: &reactiveEntity{},
		id:             id,
		name:           name,
		units:          units,
	}
	e.events = append(e.events, EquipmentCreatedEvent{
		EquipmentID:        string(e.id),
		EquipmentName:      string(e.name),
		EquipmentUnits:     string(e.units),
		EquipmentCreatedAt: time.Now().UTC(),
	})
	return e
}
