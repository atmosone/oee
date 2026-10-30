package workshifts

import "fmt"

type ID string

func (id ID) Validate() error {
	if id == "" {
		return fmt.Errorf("%w (zero id)", ErrValidation)
	}
	return nil
}

type Name string

func (n Name) Validate() error {}

type Status string

const (
	StatusActive   = Status("active")
	StatusFinished = Status("finished")
)

func (s Status) Validate() error { return nil }
