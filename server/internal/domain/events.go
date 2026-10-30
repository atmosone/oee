package domain

import (
	"time"

	"github.com/atmosone/cqrs"
)

const (
	EventTypeEquipmentLineCreated  = cqrs.EventType("equipment_line_created")
	EventTypeEquipmentLineExpanded = cqrs.EventType("equipment_line_expanded")
	EventTypeSpecificationCreated  = cqrs.EventType("specification_created")

	EventTypeDowntimeStarted          = cqrs.EventType("downtime_started")
	EventTypeDowntimeFinished         = cqrs.EventType("downtime_finished")
	EventTypeQualificationRuleCreated = cqrs.EventType("qualification_rule_created")
	EventTypeQualificationCreated     = cqrs.EventType("qualification_created")
)

type EquipmentLineExpandedEvent struct {
	EquipmentLineID         string    `json:"equipment_line_id"`
	EquipmentID             string    `json:"equipment_id"`
	EquipmentLineExpandedAt time.Time `json:"equipment_line_expanded_at"`
}

func (e EquipmentLineExpandedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentLineExpanded
}

type EquipmentCreatedEvent struct {
	EquipmentID        string    `json:"equipment_id"`
	EquipmentName      string    `json:"equipment_name"`
	EquipmentUnits     string    `json:"equipment_units"`
	EquipmentCreatedAt time.Time `json:"equipment_created_at"`
}

func (e EquipmentCreatedEvent) Type() cqrs.EventType {
	return EventTypeEquipmentCreated
}

type DowntimeStartedEvent struct {
	DowntimeID        string    `json:"downtime_id"`
	EquipmentID       string    `json:"equipment_id"`
	DowntimeStatus    string    `json:"downtime_status"`
	DowntimeStartedAt time.Time `json:"downtime_started_at"`
}

func (evt DowntimeStartedEvent) Type() cqrs.EventType {
	return EventTypeDowntimeStarted
}

type DowntimeFinishedEvent struct {
	DowntimeID         string        `json:"downtime_id"`
	EquipmentID        string        `json:"equipment_id"`
	DowntimeStatus     string        `json:"downtime_status"`
	DowntimeDuration   time.Duration `json:"downtime_duration"`
	DowntimeFinishedAt time.Time     `json:"downtime_finished_at"`
}

func (evt DowntimeFinishedEvent) Type() cqrs.EventType {
	return EventTypeDowntimeFinished
}

type SpecificationCreatedEvent struct {
	SpecificationID          string    `json:"specification_id"`
	SpecificationCode        string    `json:"specification_code"`
	SpecificationDescription string    `json:"specification_description"`
	SpecificationCreatedAt   time.Time `json:"specification_created_at"`
}

func (e SpecificationCreatedEvent) Type() cqrs.EventType {
	return EventTypeSpecificationCreated
}

type QualificationRuleCreatedEvent struct {
	EquipmentID                string    `json:"equipment_id"`
	SpecificationID            string    `json:"specification_id"`
	QualificationRuleCreatedAt time.Time `json:"qualification_rule_created_at"`
}

func (e QualificationRuleCreatedEvent) Type() cqrs.EventType {
	return EventTypeQualificationRuleCreated
}

type QualificationCreatedEvent struct {
	DowntimeID             string    `json:"downtime_id"`
	SpecificationID        string    `json:"specification_id"`
	QualificationComment   string    `json:"qualification_comment"`
	QualificationCreatedAt time.Time `json:"qualification_created_at"`
}

func (e QualificationCreatedEvent) Type() cqrs.EventType {
	return EventTypeQualificationCreated
}
