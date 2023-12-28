package assert

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// NotNil checks given value if it is not nil
func NotNil(values ...any) error {
	for i, value := range values {
		if value == nil {
			return errors.WithStack(fmt.Errorf("assert failed: given value at index %d is nil", i))
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Func:
			if reflect.ValueOf(value).IsNil() {
				return errors.WithStack(fmt.Errorf("assert failed: given value at index %d is nil", i))
			}
		}
	}

	return nil
}

// StringLength checks a string for the given min and max length
func StringLength(value string, minLength int, maxLength int, allowedValues *[]string /* can be nil */) error {
	if minLength != -1 && maxLength != -1 && minLength > maxLength {
		return errors.Errorf("assert misuse: minLength should be smaller than maxLength (minLength: %d, maxLength: %d)", minLength, maxLength)
	}

	lenString := len(value)
	if minLength != -1 && lenString < minLength {
		return errors.Errorf("assert failed: given value '%s' is too short (%d < %d)", value, lenString, minLength)
	}

	if maxLength != -1 && lenString > maxLength {
		return errors.Errorf("assert failed: given value '%s' is too long (%d > %d)", value, lenString, minLength)
	}

	if allowedValues != nil {
		found := false
		for _, allowedValue := range *allowedValues {
			if allowedValue == value {
				found = true
				break
			}
		}

		if !found {
			return errors.Errorf("assert failed: given value '%s' is not in allowed values (%s)", value, *allowedValues)
		}
	}

	return nil
}

// StringNotEmpty checks if given string is not empty
func StringNotEmpty(value string) error {
	return StringLength(value, 1, -1, nil)
}
