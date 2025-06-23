package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	val := reflect.ValueOf(v).Elem() // holds a pointer, otherwise it panics
	if val.Kind() != reflect.Struct {
		return nil // or return error if you only want to validate structs
	}
	typ := val.Type()
	fmt.Println(" .... ..... ..... ")
	fmt.Println("input:            ", v)
	fmt.Println("input val:        ", val)

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {

			fmt.Println("field: ", field)
			fmt.Println("Type of field: ", fieldType.Type)

			if fieldType.Type.Kind() == reflect.String {
				fmt.Println("Input1 is a string so that is OK ")

			} else if (fieldType.Type == reflect.TypeOf(json.RawMessage{})) {
				fmt.Println("Input2 is json.RawMessage so that is OK")

			} else {
				// Save a ValidationError indicating the field is not properly set
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldType.Name,
					Err:   fmt.Errorf("field %s is not properly set and has no validation tag", fieldType.Name),
				})
			}

			fmt.Println(".................................")
			continue
		}

		rules := strings.Split(tag, "|")
		fieldName := fieldType.Name

		fmt.Println("field value: ", field)
		fmt.Println("field Type:  ", fieldType)
		fmt.Println("pravila:     ", rules)
		fmt.Println("pole:        ", fieldName)

		for _, rule := range rules {
			parts := strings.Split(rule, ":")
			ruleName := parts[0]
			var param string
			if len(parts) > 1 {
				param = parts[1]
			}

			switch ruleName {
			case "min":
				if field.Kind() == reflect.Int {
					minVal, err := strconv.Atoi(param)
					if err != nil {
						return fmt.Errorf("invalid min parameter for field %s", fieldType.Name)
					}
					if int(field.Int()) < minVal {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("value %d less than min %d", field.Int(), minVal),
						})
					}
				}
			// Add other kinds if needed (float, etc.)
			case "max":
				if field.Kind() == reflect.Int {
					maxVal, err := strconv.Atoi(param)
					if err != nil {
						return fmt.Errorf("invalid max parameter for field %s", fieldType.Name)
					}
					if int(field.Int()) > maxVal {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("value %d greater than max %d", field.Int(), maxVal),
						})
					}
				}

			case "len":
				expectedLen, err := strconv.Atoi(param)
				if err != nil {
					return fmt.Errorf("invalid len parameter for field %s", fieldType.Name)
				}
				switch field.Kind() {
				case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
					if field.Len() != expectedLen {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("length %d does not equal expected %d", field.Len(), expectedLen),
						})
					}
				default:
					return fmt.Errorf("len validation not supported for field %s of kind %s", fieldType.Name, field.Kind())
				}

			}

		}

	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
