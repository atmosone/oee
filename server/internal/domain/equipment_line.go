package domain

import (
	"fmt"
	"time"

	"github.com/atmosone/cqrs"
)

const EventTypeEquipmentLineCreated = cqrs.EventType("equipment_line_created")
const EventTypeEquipmentLineExpanded = cqrs.EventType("equipment_line_expanded")

type EquipmentLine struct {
	*reactiveEntity
	id        EquipmentLineID
	name      EquipmentLineName
	equipment map[EquipmentID]struct{}
}

func NewEquipmentLine(id EquipmentLineID, name EquipmentLineName) *EquipmentLine {
	el := &EquipmentLine{
		reactiveEntity: &reactiveEntity{},
		id:             id,
		name:           name,
		equipment:      make(map[EquipmentID]struct{}),
	}
	el.events = append(el.events, EquipmentLineCreatedEvent{
		EquipmentLineID:        string(el.id),
		EquipmentLineName:      string(el.name),
		EquipmentLineCreatedAt: time.Now().UTC(),
	})
	return el
}

func (el *EquipmentLine) Expand(equipmentID EquipmentID) error {
	if _, ok := el.equipment[equipmentID]; ok {
		return fmt.Errorf("%w (equipment already exists in this line)", ErrConflict)
	}
	el.equipment[equipmentID] = struct{}{}
	el.events = append(el.events, EquipmentLineExpandedEvent{
		EquipmentLineID:         string(el.id),
		EquipmentID:             string(equipmentID),
		EquipmentLineExpandedAt: time.Now().UTC(),
	})
	return nil
}
