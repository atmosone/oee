package qualification

import (
	"github.com/atmosone/cqrs"
)

const (
	EventTypeRuleCreated       = cqrs.EventType("evt.qualification.rule_created")
	EventTypeRuleRemoved       = cqrs.EventType("evt.qualification.rule_removed")
	EventTypeDowntimeQualified = cqrs.EventType("evt.qualification.downtime_qualified")
)

type RuleCreatedEvent struct {
	EquipmentID     string
	SpecificationID string
}

func (e RuleCreatedEvent) Type() cqrs.EventType {
	return EventTypeRuleCreated
}

type RuleRemovedEvent struct {
	EquipmentID     string
	SpecificationID string
}

func (e RuleRemovedEvent) Type() cqrs.EventType {
	return EventTypeRuleRemoved
}

type DowntimeQualifiedEvent struct {
	DowntimeID      string
	SpecificationID string
	Comment         string
}

func (e DowntimeQualifiedEvent) Type() cqrs.EventType {
	return EventTypeDowntimeQualified
}
