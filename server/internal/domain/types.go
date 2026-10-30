package domain

import (
	"fmt"
	"time"
)

func ValidateTimestamp(time time.Time) error {
	if time.IsZero() {
		return fmt.Errorf(
			"%w (time zero value)",
			ErrValidationFailed,
		)
	}
	return nil
}

type EquipmentLineID string

func (elid EquipmentLineID) Validate() error {
	if elid == "" {
		return fmt.Errorf(
			"%w (empty equipment line id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type EquipmentLineName string

func (eln EquipmentLineName) Validate() error {
	const min, max = 1, 64
	if l := len(eln); l < min || l > max {
		return fmt.Errorf(
			"%w (equipment line name can contains from %d to %d characters)",
			ErrValidationFailed, min, max)
	}
	return nil
}

type DowntimeID string

func (did DowntimeID) Validate() error {
	if did == "" {
		return fmt.Errorf(
			"%w (empty downtime id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type DowntimeStatus string

const (
	DowntimeStatusActive   = DowntimeStatus("active")
	DowntimeStatusFinished = DowntimeStatus("finished")
)

type DowntimeCategory string

const (
	DowntimeCategoryUnknown = DowntimeCategory("unknown")
	DowntimeCategoryShort   = DowntimeCategory("short")
	DowntimeCategoryLong    = DowntimeCategory("long")
)

type QualificationRuleID string

func (qrid QualificationRuleID) Validate() error {
	if qrid == "" {
		return fmt.Errorf(
			"%w (empty qualification rule id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type SpecificationID string

func (sid SpecificationID) Validate() error {
	if sid == "" {
		return fmt.Errorf(
			"%w (empty specification id)",
			ErrValidationFailed,
		)
	}
	return nil
}

type SpecificationCode string

func (sc SpecificationCode) Validate() error {
	const min, max = 1, 25
	if l := len(sc); l < min || l > max {
		return fmt.Errorf(
			"%w (specification code can contains from %d to %d characters)",
			ErrValidationFailed, min, max)
	}
	return nil
}

type SpecificationDescription string

func (sd SpecificationDescription) Validate() error {
	const min, max = 1, 128
	if l := len(sd); l < min || l > max {
		return fmt.Errorf(
			"%w (specification description can contains from %d to %d characters)",
			ErrValidationFailed, min, max)
	}
	return nil
}

type QualificationComment string

func (qc QualificationComment) Validate() error {
	const max = 256
	if l := len(qc); l > max {
		return fmt.Errorf(
			"%w (qualification comment can contains to %d characters)",
			ErrValidationFailed, max)
	}
	return nil
}
