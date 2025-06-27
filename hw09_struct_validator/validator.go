package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("Validation Error: field %s: %v\n", ve.Field, ve.Err))
	}
	return sb.String()
}

type DeveloperError struct {
	Msg string
}

func (e DeveloperError) Error() string {
	return fmt.Sprintf("Developer error: %s", e.Msg)
}

var devErr DeveloperError

// var valErrs ValidationErrors // TO-Do remove as the switch is using default for this one as well

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	val, err := validateInput(v)
	if err != nil {
		// Developer error — return immediately
		return err
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		fieldName := fieldType.Name

		if tag == "" {
			continue // absence of tag is not an error
		}

		rules := strings.Split(tag, "|")

		for _, rule := range rules {
			parts := strings.SplitN(rule, ":", 2)
			ruleName := parts[0]
			var param string
			if len(parts) > 1 {
				param = parts[1]
			}

			var err error
			switch ruleName {
			case "min":
				err = validateMinRule(field, param, fieldName)
			case "max":
				err = validateMaxRule(field, param, fieldName)
			case "len":
				err = validateLenRule(field, param, fieldName)
			case "regexp":
				err = validateRegexpRule(field, param, fieldName)
			case "in":
				err = validateInRule(field, param, fieldName)
			default:
				err = DeveloperError{Msg: fmt.Sprintf("unknown validation rule %s for field %s", ruleName, fieldName)}
			}

			if err != nil {
				switch {
				case errors.As(err, &devErr):
					// Developer error — stop and return immediately
					return err
				default:
					// Validation error or other errors — collect and continue
					validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateInput(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return reflect.Value{}, DeveloperError{Msg: "nil pointer passed to validator"}
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return reflect.Value{}, DeveloperError{Msg: fmt.Sprintf("not a struct type passed to validator, got %s", val.Kind())}
	}

	return val, nil
}

func validateMinRule(field reflect.Value, param string, fieldName string) error {
	if !isIntKind(field.Kind()) {
		return DeveloperError{Msg: fmt.Sprintf("min validation only supports int kinds, field %s is %s",
			fieldName, field.Kind())}
	}
	minVal, err := strconv.Atoi(param)
	if err != nil {
		return DeveloperError{Msg: fmt.Sprintf("invalid min parameter for field %s", fieldName)}
	}
	if int(field.Int()) < minVal {
		return fmt.Errorf("value %d less than min %d", field.Int(), minVal) // validation error
	}
	return nil
}

func validateMaxRule(field reflect.Value, param string, fieldName string) error {
	if !isIntKind(field.Kind()) {
		return DeveloperError{Msg: fmt.Sprintf("max validation only supports int kinds, field %s is %s",
			fieldName, field.Kind())}
	}
	maxVal, err := strconv.Atoi(param)
	if err != nil {
		return DeveloperError{Msg: fmt.Sprintf("invalid max parameter for field %s", fieldName)}
	}
	if int(field.Int()) > maxVal {
		return fmt.Errorf("value %d greater than max %d", field.Int(), maxVal) // validation error
	}
	return nil
}

func validateLenRule(field reflect.Value, param string, fieldName string) error {
	expectedLen, err := strconv.Atoi(param)
	if err != nil {
		return DeveloperError{Msg: fmt.Sprintf("invalid len parameter for field %s", fieldName)}
	}
	switch {
	case field.Kind() == reflect.String:
		if field.Len() != expectedLen {
			return fmt.Errorf("length %d does not equal expected %d", field.Len(), expectedLen)
		}
	case field.Kind() == reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)
			if elem.Kind() != reflect.String {
				return DeveloperError{Msg: fmt.Sprintf("len validation: element %d in field %s is not a string", i, fieldName)}
			}
			if elem.Len() != expectedLen {
				return fmt.Errorf("element %d length %d does not equal expected %d", i, elem.Len(), expectedLen)
			}
		}
	default:
		return DeveloperError{Msg: fmt.Sprintf("len validation not supported for field %s of kind %s",
			fieldName, field.Kind())}
	}
	return nil
}

func validateRegexpRule(field reflect.Value, param string, fieldName string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return DeveloperError{Msg: fmt.Sprintf("invalid regexp pattern for field %s: %v", fieldName, err)}
	}
	if field.Kind() != reflect.String {
		return DeveloperError{Msg: fmt.Sprintf("regexp validation only for string, field %s is %s", fieldName, field.Kind())}
	}
	if !re.MatchString(field.String()) {
		return fmt.Errorf("field %s does not match regexp %s", fieldName, param) // validation error
	}
	return nil
}

func validateInRule(field reflect.Value, param string, fieldName string) error {
	allowed := strings.Split(param, ",")
	found := false

	switch {
	case field.Kind() == reflect.String:
		val := field.String()
		for _, a := range allowed {
			if val == a {
				found = true
				break
			}
		}
	case isIntKind(field.Kind()):
		var allowedInts []int64
		for _, s := range allowed {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return DeveloperError{Msg: fmt.Sprintf("invalid integer value %q in 'in' rule: %v", s, err)}
			}
			allowedInts = append(allowedInts, n)
		}
		val := field.Int()
		for _, a := range allowedInts {
			if val == a {
				found = true
				break
			}
		}
	default:
		return DeveloperError{Msg: fmt.Sprintf("in validation only supports string & int kinds, field %s is %s",
			fieldName, field.Kind())}
	}

	if !found {
		return fmt.Errorf("value %q not allowed, must be one of %v", field.Interface(), allowed) // validation error
	}
	return nil
}

func isIntKind(k reflect.Kind) bool {
	switch {
	case k == reflect.Int:
		// TO-DO: reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}
