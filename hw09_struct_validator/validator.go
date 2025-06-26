package hw09structvalidator

import (
	"encoding/json"
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
		sb.WriteString(fmt.Sprintf("Field %s: %v\n", ve.Field, ve.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	val, err := validateInput(v)
	if err != nil {
		return err
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		fieldName := fieldType.Name

		if tag == "" {
			switch {
			case fieldType.Type.Kind() == reflect.String:
				if field.String() == "" {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("field %s empty and has no validation tag", fieldName),
					})
				}
			case fieldType.Type == reflect.TypeOf(json.RawMessage{}):
				// No validation for json.RawMessage without tags
			case fieldType.Type.Kind() == reflect.Slice && fieldType.Type.Elem().Kind() == reflect.Uint8:
				// []byte field
				if field.Len() == 0 {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   fmt.Errorf("field %s is empty and has no validation tag", fieldName),
					})
				}
			default:
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   fmt.Errorf("field %s is not properly set and has no validation tag", fieldName),
				})
			}
			continue
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
				err = fmt.Errorf("unknown validation rule %s for field %s", ruleName, fieldName)
			}

			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
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
			return reflect.Value{}, ValidationErrors{
				ValidationError{
					Field: "",
					Err:   fmt.Errorf("input error - nil pointer passed to validator"),
				},
			}
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return reflect.Value{}, ValidationErrors{
			ValidationError{
				Field: "",
				Err:   fmt.Errorf("input error - not a struct type passed to validator, got %s", val.Kind()),
			},
		}
	}

	return val, nil
}

func validateMinRule(field reflect.Value, param string, fieldName string) error {
	if !isIntKind(field.Kind()) {
		return fmt.Errorf("min validation only supports int kinds, field %s is %s", fieldName, field.Kind())
	}
	minVal, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid min parameter for field %s", fieldName)
	}
	if int(field.Int()) < minVal {
		return fmt.Errorf("value %d less than min %d", field.Int(), minVal)
	}
	return nil
}

func validateMaxRule(field reflect.Value, param string, fieldName string) error {
	if !isIntKind(field.Kind()) {
		return fmt.Errorf("max validation only supports int kinds, field %s is %s", fieldName, field.Kind())
	}
	maxVal, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid max parameter for field %s", fieldName)
	}
	if int(field.Int()) > maxVal {
		return fmt.Errorf("value %d greater than max %d", field.Int(), maxVal)
	}
	return nil
}

func validateLenRule(field reflect.Value, param string, fieldName string) error {
	expectedLen, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid len parameter for field %s", fieldName)
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
				return fmt.Errorf("len validation: element %d in field %s is not a string", i, fieldName)
			}
			if elem.Len() != expectedLen {
				return fmt.Errorf("element %d length %d does not equal expected %d", i, elem.Len(), expectedLen)
			}
		}
	default:
		return fmt.Errorf("len validation not supported for field %s of kind %s", fieldName, field.Kind())
	}
	return nil
}

func validateRegexpRule(field reflect.Value, param string, fieldName string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return fmt.Errorf("invalid regexp pattern for field %s: %w", fieldName, err)
	}
	if field.Kind() != reflect.String {
		return fmt.Errorf("regexp validation only for string, field %s is %s", fieldName, field.Kind())
	}
	if !re.MatchString(field.String()) {
		return fmt.Errorf("field %s does not match regexp %s", fieldName, param)
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
	case field.Kind() == reflect.Int:
		// TO-DO: reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var allowedInts []int64
		for _, s := range allowed {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value %q in 'in' rule: %w", s, err)
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
		return fmt.Errorf("in validation only supports string & int kinds, field %s is %s", fieldName, field.Kind())
	}

	if !found {
		return fmt.Errorf("value %q not allowed, must be one of %v", field.Interface(), allowed)
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
