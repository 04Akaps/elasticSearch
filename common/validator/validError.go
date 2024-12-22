package validator

import "fmt"

type ErrValidation struct {
	FailedField string
	Tag         string
	Value       interface{}
}

func NewErrValidation(failedField string, tag string, value interface{}) *ErrValidation {
	return &ErrValidation{FailedField: failedField, Tag: tag, Value: value}
}

func (v *ErrValidation) Error() string {
	return fmt.Sprintf(
		"[%s]: '%v' | Needs to implement '%s'",
		v.FailedField,
		v.Value,
		v.Tag,
	)
}
