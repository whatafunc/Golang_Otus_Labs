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
		// Valid User
		{
			name: "Valid User",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Alice",
				Age:    30,
				Email:  "alice@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			expectValid: true,
		},
		// Valid App
		{
			name: "Valid App",
			in: &App{
				Version: "1.3.5",
			},
			expectValid: true,
		},
		// Valid Response
		{
			name: "Valid Response",
			in: &Response{
				Code: 200,
				Body: "OK",
			},
			expectValid: true,
		},
		// Token struct (no validation tags)
		{
			name: "Token no tags",
			in: &Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectValid: true,
		},
		// Developer error: invalid tag param (len:xvii)
		{
			name: "Developer error invalid len param",
			in: struct {
				Field string `validate:"len:xvii"`
			}{
				Field: "test",
			},
			expectDev: true,
		},
		// Developer error: regex on int
		{
			name: "Developer error regex on int",
			in: struct {
				Field int `validate:"regexp:^\\d+$"`
			}{
				Field: 123,
			},
			expectDev: true,
		},
		// Developer error: unknown rule
		{
			name: "Developer error unknown rule",
			in: struct {
				Field string `validate:"unknown:123"`
			}{
				Field: "abc",
			},
			expectDev: true,
		},
		// Developer error: invalid in param
		{
			name: "Developer error invalid in param",
			in: struct {
				Field int `validate:"in:1,2,a"`
			}{
				Field: 1,
			},
			expectDev: true,
		},
		// Developer error: passing non-struct
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

			// For other cases (invalid user input), expect ValidationErrors
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
		// Invalid User: ID too short
		{
			name: "Invalid User ID too short",
			in: &User{
				ID:     "12345678-1",
				Name:   "Vasya",
				Age:    45,
				Email:  "Vasya@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
		},
		// Invalid User: Age too low
		{
			name: "Invalid User Age too low",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Carol",
				Age:    15,
				Email:  "carol@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
		},
		// Invalid User: wrong phone!
		{
			name: "Invalid User wrong phone",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Caroline",
				Age:    31,
				Email:  "caroline@mail.com",
				Role:   "stuff",
				Phones: []string{"02", "112"},
			},
		},
		// Invalid User: Email invalid
		{
			name: "Invalid User Email invalid",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Dave",
				Age:    32,
				Email:  "not-an-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
		},
		// Invalid User: Role not allowed
		{
			name: "Invalid User Role not allowed",
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Eve",
				Age:    30,
				Email:  "eve@mail.com",
				Role:   "user",
				Phones: []string{"12345678901"},
			},
		},
		// Invalid User: all wrong fields
		{
			name: "Invalid User all wrong fields",
			in: &User{
				ID:     "12",
				Name:   "",
				Age:    310,
				Email:  "eve!mail.com",
				Role:   "abuser",
				Phones: []string{""},
			},
		},
		// Invalid App: Version too short
		{
			name: "Invalid App Version too short",
			in: &App{
				Version: "1.0",
			},
		},
		// Invalid Response: Code not allowed
		{
			name: "Invalid Response Code not allowed",
			in: &Response{
				Code: 9999,
				Body: "Created",
			},
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
