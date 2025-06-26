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

	val := reflect.ValueOf(v)

	// Check for nil pointer before Elem()
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			validationErrors = append(validationErrors, ValidationError{
				Field: "",
				Err:   fmt.Errorf("input error - nil pointer passed to validator"),
			})
			return validationErrors
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		validationErrors = append(validationErrors, ValidationError{
			Field: "",
			Err:   fmt.Errorf("input error - not a struct type passed to validator, got %s", val.Kind()),
		})
		return validationErrors
	}

	// val := reflect.ValueOf(v).Elem() // holds a pointer, otherwise it panics
	// if val.Kind() == reflect.Ptr {
	// 	val = val.Elem()
	// }
	// if val.Kind() != reflect.Struct {
	// 	// If v is a pointer, get the element it points to

	// 	validationErrors = append(validationErrors, ValidationError{
	// 		Field: "",
	// 		Err:   fmt.Errorf("Input error - not a strcut type passed to validator %s", val.Kind()),
	// 	})
	// }
	typ := val.Type()
	fmt.Println(" .... ..... ..... ")
	fmt.Println("input:            ", v)
	fmt.Println("input val:        ", val)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			fmt.Println("field: ", field)
			fmt.Println("Type of field: ", fieldType.Type)

			switch {
			case fieldType.Type.Kind() == reflect.String:
				fmt.Println("Input1 is a string so that is OK ")
				if field.String() == "" {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   fmt.Errorf("field %s empty and has no validation tag", fieldType.Name),
					})
				}
			case fieldType.Type == reflect.TypeOf(json.RawMessage{}):
				fmt.Println("Input2 is json.RawMessage so that is OK")

			case fieldType.Type.Kind() == reflect.Slice && fieldType.Type.Elem().Kind() == reflect.Uint8:
				// []byte field
				if field.Len() == 0 {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   fmt.Errorf("field %s is empty and has no validation tag", fieldType.Name),
					})
				}

			default:
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
				case reflect.String, reflect.Array, reflect.Map:
					if field.Len() != expectedLen {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("length %d does not equal expected %d", field.Len(), expectedLen),
						})
					}
				case reflect.Slice:
					// For slice, check each element if it's a string
					for i := 0; i < field.Len(); i++ {
						elem := field.Index(i)
						if elem.Kind() != reflect.String {
							return fmt.Errorf("len validation: element %d in field %s is not a string", i, fieldType.Name)
						}
						if elem.Len() != expectedLen {
							validationErrors = append(validationErrors, ValidationError{
								Field: fieldType.Name,
								Err:   fmt.Errorf("element %d length %d does not equal expected %d", i, elem.Len(), expectedLen),
							})
						}
					}
				// fixing go lint issues:
				case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
					reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
					reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr, reflect.Struct, reflect.UnsafePointer:
					return fmt.Errorf("len validation not supported for field %s of kind %s", fieldType.Name, field.Kind())
				}
			case "regexp":
				pattern := strings.TrimPrefix(rule, "regexp:")

				// Compile regexp once per field per validation call
				re, err := regexp.Compile(pattern)
				if err != nil {
					return fmt.Errorf("invalid regexp pattern for field %s: %w", fieldType.Name, err)
				}

				// Only apply regexp if field is string
				if field.Kind() == reflect.String {
					if !re.MatchString(field.String()) {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("field %s does not match regexp %s", fieldType.Name, pattern),
						})
					}
				} else {
					return fmt.Errorf("regexp validation only for string, field %s is %s", fieldType.Name, field.Kind())
				}
			case "in":
				allowed := strings.Split(strings.TrimPrefix(rule, "in:"), ",")
				var fieldVal interface{}
				found := false
				if field.Kind() != reflect.String {
					if field.Kind() != reflect.Int {
						validationErrors = append(validationErrors, ValidationError{
							Field: fieldType.Name,
							Err:   fmt.Errorf("in validation only supports string & int kinds, field %s is %s", fieldType.Name, field.Kind()),
						})
					} else {

						var allowedInts []int
						for _, s := range allowed {
							n, err := strconv.Atoi(s)
							if err != nil {
								return fmt.Errorf("invalid integer value %q in 'in' rule: %v", s, err)
							}
							allowedInts = append(allowedInts, n)
						}

						fieldVal = int(field.Int())

						for _, intVal := range allowedInts {
							if fieldVal == intVal {
								found = true
								break
							}
						}

						fmt.Println(" ahtung ====== ")
						fmt.Println(" value: ", fieldVal)
						fmt.Println(" ahtung end====== ")
					}
				} else {
					fieldVal = field.String()
					for _, a := range allowed {
						if fieldVal == a {
							found = true
							break
						}
					}
				}

				if !found {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   fmt.Errorf("value %q not allowed, must be one of %v", fieldVal, allowed),
					})
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
