package assert

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

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

// ValidProjectID check if given string is a valid project ID
func ValidProjectID(value string) error {
	if err := StringNotEmpty(value); err != nil {
		return err
	}

	if len(value) < 5 {
		return errors.Errorf("assert failed: given value '%s' is too short", value)
	}

	if value[:4] != "pro-" {
		return errors.Errorf("assert failed: given value '%s' does not start with 'pro-'", value)
	}

	return nil
}

// ValidAPISecret checks if given string isa  valid API secret
func ValidAPISecret(value string) error {
	if err := StringNotEmpty(value); err != nil {
		return err
	}

	if len(value) < 10 {
		return errors.Errorf("assert failed: given value '%s' is too short", value)
	}

	if value[:9] != "corbado1_" {
		return errors.Errorf("assert failed: given value '%s' does not start with 'corbado1_'", value)
	}

	return nil
}

// ValidAPIEndpoint check if given string is a valid API endpoint
func ValidAPIEndpoint(value string) error {
	if err := StringNotEmpty(value); err != nil {
		return err
	}

	u, err := url.Parse(value)
	if err != nil {
		return errors.WithStack(err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return errors.Errorf("assert failed: scheme needs to be 'https' or 'http' in given value '%s' (scheme: '%s')", value, u.Scheme)
	}

	if u.Host == "" {
		return errors.Errorf("assert failed: host must not be empty in given value '%s'", value)
	}

	if u.User.Username() != "" {
		return errors.Errorf("assert failed: username must be empty in given value '%s' (username: '%s')", value, u.User.Username())
	}

	password, _ := u.User.Password()
	if password != "" {
		return errors.Errorf("assert failed: password must be empty in given value '%s' (password: '%s')", value, password)
	}

	if u.Path != "" {
		return errors.Errorf("assert failed: path must be empty in given value '%s' (path: '%s')", value, u.Path)
	}

	if u.Fragment != "" {
		return errors.Errorf("assert failed: fragment must be empty in given value '%s' (fragment: '%s')", value, u.Fragment)
	}

	if u.RawQuery != "" {
		return errors.Errorf("assert failed: querystring must be empty in given value '%s' (querystring: '%s')", value, u.RawQuery)
	}

	return nil
}

// DurationNotEmpty checks if given duration is not empty
func DurationNotEmpty(value time.Duration) error {
	if value == 0 {
		return errors.Errorf("assert failed: given value '%s' is empty", value)
	}

	return nil
}
