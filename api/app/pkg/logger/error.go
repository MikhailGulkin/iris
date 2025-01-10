package log

import (
	"errors"
	"fmt"
)

type errWithFields interface {
	Fields() Fld
	Origin() error
}

type FieldsError struct {
	err    error
	fields Fld
}

func (e *FieldsError) Error() string {
	return e.err.Error()
}

func (e *FieldsError) Is(target error) bool {
	return errors.Is(e.err, target)
}

func (e *FieldsError) Fields() Fld {
	return e.fields
}

func (e *FieldsError) Origin() error {
	return e.err
}

// Wrap error with fields for logging
func Wrap(msg string, err error, fields Fld) error {
	if fields == nil {
		fields = Fld{}
	}
	fieldsErr, ok := err.(errWithFields)
	if !ok {
		return &FieldsError{
			err:    fmt.Errorf("%s: %w", msg, err),
			fields: fields,
		}
	}

	return &FieldsError{
		err:    fmt.Errorf("%s: %w", msg, fieldsErr.Origin()),
		fields: mergeFields(fieldsErr.Fields(), fields),
	}
}

func mergeFields(fld1, fld2 Fld) Fld {
	result := fld1
	for k, v := range fld2 {
		result[k] = v
	}

	return result
}
