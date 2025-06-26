package hw09structvalidator

import (
	"encoding/json"
	"fmt"
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

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Valid User
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Alice",
				Age:    30,
				Email:  "alice@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			expectedErr: nil,
		},
		// Valid App
		{
			in: &App{
				Version: "1.3.5",
			},
			expectedErr: nil,
		},
		// Valid Response
		{
			in: &Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		// Token struct (no validation tags)
		{
			in: &Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			testResults := Validate(tt.in)
			if (tt.expectedErr != nil && testResults == nil) || (tt.expectedErr == nil && testResults != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, testResults)
			} else {
				fmt.Println("----VALIDATION TESTS PASSED-----")
				if testResults != nil {
					// Safe to call Error() since testResults is not nil
					fmt.Println("Validation caught correct assumptions aka errors:\n", testResults.Error())
				} else {
					fmt.Println("Validation passed with no errors")
				}
			}
			_ = tt
		})
	}
}

func TestInvalidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Invalid User: ID too short
		{
			in: &User{
				ID:     "12345678-1",
				Name:   "Vasya",
				Age:    45,
				Email:  "Vasya@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid User: Age too low
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Carol",
				Age:    15,
				Email:  "carol@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid User: wrong phone
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Caroline",
				Age:    31,
				Email:  "caroline@mail.com",
				Role:   "stuff",
				Phones: []string{"02", "112"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid User: Email invalid
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Dave",
				Age:    32,
				Email:  "not-an-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid User: Role not allowed
		{
			in: &User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Eve",
				Age:    30,
				Email:  "eve@mail.com",
				Role:   "user",
				Phones: []string{"12345678901"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid User: all wrong fields
		{
			in: &User{
				ID:     "12",
				Name:   "",
				Age:    310,
				Email:  "eve!mail.com",
				Role:   "abuser",
				Phones: []string{""},
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid App: Version too short
		{
			in: &App{
				Version: "1.0",
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Invalid Response: Code not allowed
		{
			in: &Response{
				Code: 9999,
				Body: "Created",
			},
			expectedErr: fmt.Errorf("validation error"),
		},
		// Non-struct type (int)
		{
			in:          42,
			expectedErr: fmt.Errorf("validation error"),
		},
		// Non-struct type (string)
		{
			in:          "not a struct",
			expectedErr: fmt.Errorf("validation error"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			testResults := Validate(tt.in)
			if (tt.expectedErr != nil && testResults == nil) || (tt.expectedErr == nil && testResults != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, testResults)
			} else {
				fmt.Println("----INVALIDATION TESTS PASSED-----")
				if testResults != nil {
					// Safe to call Error() since testResults is not nil
					fmt.Println("Invalidation caught correct assumptions aka errors:\n", testResults.Error())
				} else {
					fmt.Println("Invalidation passed with no errors")
				}
			}
			_ = tt
		})
	}
}
