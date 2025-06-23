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
		Phones []string        `validate:"len:2"`
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
		// Invalid User: ID too short
		{
			in: &User{
				ID:     "12345678-1",
				Name:   "Carol",
				Age:    45,
				Email:  "carol1@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: fmt.Errorf("validation error"),
		},

		// Invalid User: Name not a string
		//

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
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			testResults := Validate(tt.in)
			if (tt.expectedErr != nil && testResults == nil) || (tt.expectedErr == nil && testResults != nil) {

				// fmt.Println("Validation failed:", err)
				// err.Error() //   print all validation errors concatenated

				t.Errorf("expected error: %v, got: %v", tt.expectedErr, testResults)
			} else {
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
