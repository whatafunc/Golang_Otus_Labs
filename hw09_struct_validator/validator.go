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

type compareType int

const (
	compareMin compareType = iota
	compareMax
)

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
		return reflect.Value{}, DeveloperError{Msg: fmt.Sprintf("not a struct type passed to validator, got %s",
			val.Kind())}
	}

	return val, nil
}

func validateIntBoundaryRule(field reflect.Value, param string, fieldName string, cmpType compareType) error {
	limit, err := strconv.Atoi(param)
	if err != nil {
		return DeveloperError{Msg: fmt.Sprintf("invalid parameter for field %s", fieldName)}
	}

	check := func(val int) bool {
		switch cmpType {
		case compareMin:
			return val < limit
		case compareMax:
			return val > limit
		default:
			return false
		}
	}

	var errMsgSingle, errMsgElement string
	switch cmpType {
	case compareMin:
		errMsgSingle = "value %d less than min %d"
		errMsgElement = "element %d value %d less than min %d"
	case compareMax:
		errMsgSingle = "value %d greater than max %d"
		errMsgElement = "element %d value %d greater than max %d"
	}

	switch {
	case field.Kind() == reflect.Int:
		val := int(field.Int())
		if check(val) {
			return fmt.Errorf(errMsgSingle, val, limit)
		}
	case field.Kind() == reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)
			if !isIntKind(elem.Kind()) {
				return DeveloperError{Msg: fmt.Sprintf("validation only supports int kinds, element %d of field %s is %s",
					i, fieldName, elem.Kind())}
			}
			val := int(elem.Int())
			if check(val) {
				return fmt.Errorf(errMsgElement, i, val, limit)
			}
		}
	default:
		return DeveloperError{Msg: fmt.Sprintf("validation only supports int kinds or slices of int, field %s is %s",
			fieldName, field.Kind())}
	}

	return nil
}

func validateMinRule(field reflect.Value, param string, fieldName string) error {
	return validateIntBoundaryRule(field, param, fieldName, compareMin)
}

func validateMaxRule(field reflect.Value, param string, fieldName string) error {
	return validateIntBoundaryRule(field, param, fieldName, compareMax)
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

	switch {
	case field.Kind() == reflect.String:
		if !re.MatchString(field.String()) {
			return fmt.Errorf("field %s does not match regexp %s", fieldName, param) // validation error
		}
	case field.Kind() == reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)
			if elem.Kind() != reflect.String {
				return DeveloperError{
					Msg: fmt.Sprintf("regexp validation only supports string slices, element %d of field %s is %s",
						i, fieldName, elem.Kind()),
				}
			}
			if !re.MatchString(elem.String()) {
				return fmt.Errorf("element %d of field %s does not match regexp %s", i, fieldName, param) // validation error
			}
		}
	default:
		return DeveloperError{Msg: fmt.Sprintf("regexp validation only supports string or []string, field %s is %s",
			fieldName, field.Kind())}
	}

	return nil
}

func validateInRule(field reflect.Value, param string, fieldName string) error {
	allowedStr := strings.Split(param, ",")
	for i := range allowedStr {
		allowedStr[i] = strings.TrimSpace(allowedStr[i])
	}

	switch {
	case field.Kind() == reflect.String:
		return validateInSingleString(field.String(), allowedStr)
	case field.Kind() == reflect.Int:
		allowedInts, err := parseAllowedInts(allowedStr)
		if err != nil {
			return err
		}
		return validateInSingleInt(field.Int(), allowedInts)
	case field.Kind() == reflect.Slice:
		if field.Len() == 0 {
			return nil // empty slice passes
		}
		elemKind := field.Index(0).Kind()
		switch {
		case elemKind == reflect.String:
			return validateInSliceString(field, allowedStr, fieldName)
		case elemKind == reflect.Int:
			allowedInts, err := parseAllowedInts(allowedStr)
			if err != nil {
				return err
			}
			return validateInSliceInt(field, allowedInts, fieldName)
		default:
			return DeveloperError{
				Msg: fmt.Sprintf("in validation only supports string & int kinds, field %s element kind %s",
					fieldName, elemKind),
			}
		}
	default:
		return DeveloperError{
			Msg: fmt.Sprintf("in validation only supports string & int kinds or slices thereof, field %s is %s",
				fieldName, field.Kind()),
		}
	}
}

func parseAllowedInts(allowedStr []string) ([]int64, error) {
	allowedInts := make([]int64, 0, len(allowedStr))

	for _, s := range allowedStr {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, DeveloperError{Msg: fmt.Sprintf("invalid integer value %q in 'in' rule: %v", s, err)}
		}
		allowedInts = append(allowedInts, n)
	}
	return allowedInts, nil
}

func validateInSingleString(val string, allowed []string) error {
	for _, a := range allowed {
		if val == a {
			return nil
		}
	}
	return fmt.Errorf("value %q not allowed, must be one of %v", val, allowed)
}

func validateInSingleInt(val int64, allowed []int64) error {
	for _, a := range allowed {
		if val == a {
			return nil
		}
	}
	return fmt.Errorf("value %d not allowed, must be one of %v", val, allowed)
}

func validateInSliceString(field reflect.Value, allowed []string, fieldName string) error {
	for i := 0; i < field.Len(); i++ {
		val := field.Index(i).String()
		if err := validateInSingleString(val, allowed); err != nil {
			return fmt.Errorf("element %d of %s has err: %w", i, fieldName, err)
		}
	}
	return nil
}

func validateInSliceInt(field reflect.Value, allowed []int64, fieldName string) error {
	for i := 0; i < field.Len(); i++ {
		val := field.Index(i).Int()
		if err := validateInSingleInt(val, allowed); err != nil {
			return fmt.Errorf("element %d %s has err: %w", i, fieldName, err)
		}
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
