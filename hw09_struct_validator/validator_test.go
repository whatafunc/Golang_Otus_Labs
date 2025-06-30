package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Family struct {
		ChildrenAges  []int    `validate:"min:0|max:18"`
		ParentsStatus []string `validate:"regexp:^(Mr\.|Mrs\.)\s+[A-Za-z]+$"` //nolint
		Country       []string `validate:"in:Spain,UAE,GB"`
	}

	FamilyWithWrongStatusDetect struct {
		ChildrenAges []int `validate:"min:0|max:18"`
		//nolint:structcheck
		ParentsStatus []string `validate:"regexp:^(YYY\\.|XXX\\.)\\s+[A-Za-z]+("`
	}

	Scores struct {
		Values []int `validate:"in:1,2,3,10,23"`
	}
)

var (
	// Reusable test fixtures.
	validUser = &User{
		ID:     "12345678-1234-1234-1234-123456789012",
		Name:   "Alice",
		Age:    30,
		Email:  "alice@mail.com",
		Role:   "admin",
		Phones: []string{"12345678901", "10987654321"},
	}

	validFamily = &Family{
		ChildrenAges:  []int{5, 18},
		ParentsStatus: []string{"Mrs. Bulkina", "Mr. Ivanov"},
	}

	validApp = &App{
		Version: "1.3.5",
	}

	validResponse = &Response{
		Code: 200,
		Body: "OK",
	}

	tokenNoTags = &Token{
		Header:    []byte("header"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}

	validScores = &Scores{
		Values: []int{1, 2, 3},
	}

	// Reusable invalid test fixtures.
	invalidUserIDTooShort = &User{
		ID:     "12345678-1",
		Name:   "Vasya",
		Age:    45,
		Email:  "Vasya@mail.com",
		Role:   "admin",
		Phones: []string{"12345678901"},
	}

	invalidOldFamily = &Family{
		ChildrenAges:  []int{-1, 31},
		ParentsStatus: []string{"Cheburashka", "DOOM III", "Go Professional"},
		Country:       []string{"Sputnik", "Mars"},
	}

	invalidFamilyWithWrongStatus = &FamilyWithWrongStatusDetect{
		ChildrenAges:  []int{5, 10},
		ParentsStatus: []string{"Mrs. Bulkina", "Mr. Ivanov"},
	}

	invalidUserAgeTooLow = &User{
		ID:     "12345678-1234-1234-1234-123456789012",
		Name:   "Carol",
		Age:    15,
		Email:  "carol@mail.com",
		Role:   "admin",
		Phones: []string{"12345678901"},
	}

	invalidUserWrongPhone = &User{
		ID:     "12345678-1234-1234-1234-123456789012",
		Name:   "Caroline",
		Age:    31,
		Email:  "caroline@mail.com",
		Role:   "stuff",
		Phones: []string{"02", "112"},
	}

	invalidUserEmailInvalid = &User{
		ID:     "12345678-1234-1234-1234-123456789012",
		Name:   "Dave",
		Age:    32,
		Email:  "not-an-email",
		Role:   "admin",
		Phones: []string{"12345678901"},
	}

	invalidUserRoleNotAllowed = &User{
		ID:     "12345678-1234-1234-1234-123456789012",
		Name:   "Eve",
		Age:    30,
		Email:  "eve@mail.com",
		Role:   "user",
		Phones: []string{"12345678901"},
	}

	invalidUserAllWrongFields = &User{
		ID:     "12",
		Name:   "",
		Age:    310,
		Email:  "eve!mail.com",
		Role:   "abuser",
		Phones: []string{""},
	}

	invalidAppVersionTooShort = &App{
		Version: "1.0",
	}

	invalidResponseCodeNotAllowed = &Response{
		Code: 9999,
		Body: "Created",
	}

	invalidScoresDisallowedValues = &Scores{
		Values: []int{1, 55555, 3000},
	}
)

func isValidationErrors(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}

func isDeveloperError(err error) bool {
	var de DeveloperError
	return errors.As(err, &de)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectValid bool // true if expect no validation errors
		expectDev   bool // true if expect developer error
	}{
		{
			name:        "Valid User",
			in:          validUser,
			expectValid: true,
		},
		{
			name:        "Still a valid Family struct use",
			in:          validFamily,
			expectValid: true,
		},
		{
			name:        "Valid App",
			in:          validApp,
			expectValid: true,
		},
		{
			name:        "Valid Response",
			in:          validResponse,
			expectValid: true,
		},
		{
			name:        "Token no tags",
			in:          tokenNoTags,
			expectValid: true,
		},
		// Developer error cases
		{
			name: "Developer error invalid len param",
			in: struct {
				Field string `validate:"len:xvii"`
			}{
				Field: "test",
			},
			expectDev: true,
		},
		{
			name: "Developer error regex on int",
			in: struct {
				Field int `validate:"regexp:^\\d+$"`
			}{
				Field: 123,
			},
			expectDev: true,
		},
		{
			name: "Developer error unknown rule",
			in: struct {
				Field string `validate:"unknown:123"`
			}{
				Field: "abc",
			},
			expectDev: true,
		},
		{
			name: "Developer error invalid in param",
			in: struct {
				Field int `validate:"in:1,2,a"`
			}{
				Field: 1,
			},
			expectDev: true,
		},
		{
			name:      "Developer error non-struct int",
			in:        42,
			expectDev: true,
		},
		{
			name:      "Developer error non-struct string",
			in:        "not a struct",
			expectDev: true,
		},
		{
			name:        "Valid Scores slice",
			in:          validScores,
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)

			if tt.expectDev {
				if !isDeveloperError(err) {
					t.Errorf("expected developer error, got: %v", err)
				} else {
					t.Logf("correctly got developer error: %v", err)
				}
				return
			}

			if tt.expectValid {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				} else {
					t.Log("validation passed with no errors")
				}
				return
			}

			if err == nil {
				t.Errorf("expected validation errors, got nil")
				return
			}
			if !isValidationErrors(err) {
				t.Errorf("expected validation errors, got different error: %v", err)
				return
			}
			t.Logf("validation errors as expected:\n%s", err.Error())
		})
	}
}

func TestInvalidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectDev   bool
		expectValid bool
	}{
		{
			name:        "Invalid User ID too short",
			in:          invalidUserIDTooShort,
			expectValid: false,
		},
		{
			name:        "Old family, incorrect parents, not applicable country",
			in:          invalidOldFamily,
			expectValid: false,
		},
		{
			name:        "Not OK family",
			in:          invalidFamilyWithWrongStatus,
			expectValid: false,
			expectDev:   true,
		},
		{
			name:        "Invalid User Age too low",
			in:          invalidUserAgeTooLow,
			expectValid: false,
		},
		{
			name:        "Invalid User wrong phone",
			in:          invalidUserWrongPhone,
			expectValid: false,
		},
		{
			name:        "Invalid User Email invalid",
			in:          invalidUserEmailInvalid,
			expectValid: false,
		},
		{
			name:        "Invalid User Role not allowed",
			in:          invalidUserRoleNotAllowed,
			expectValid: false,
		},
		{
			name:        "Invalid User all wrong fields",
			in:          invalidUserAllWrongFields,
			expectValid: false,
		},
		{
			name:        "Invalid App Version too short",
			in:          invalidAppVersionTooShort,
			expectValid: false,
		},
		{
			name:        "Invalid Response Code not allowed",
			in:          invalidResponseCodeNotAllowed,
			expectValid: false,
		},
		{
			name:        "Invalid Scores slice with disallowed value",
			in:          invalidScoresDisallowedValues,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tt.in)

			if tt.expectDev {
				if !isDeveloperError(err) {
					t.Errorf("expected developer error, got: %v", err)
				} else {
					t.Logf("correctly got developer error: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("expected validation errors, got nil")
				return
			}
			if !isValidationErrors(err) {
				t.Errorf("expected validation errors, got different error: %v", err)
				return
			}
			t.Logf("validation errors as expected:\n%s", err.Error())
		})
	}
}
