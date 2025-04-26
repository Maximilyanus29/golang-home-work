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

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "98765432109"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user ID length",
			in: User{
				ID:     "short-id",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrInvalidValidation},
			},
		},
		{
			name: "invalid user age (too young)",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    16,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrInvalidValidation},
			},
		},
		{
			name: "invalid user email format",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: ErrInvalidValidation},
			},
		},
		{
			name: "invalid user role",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "user",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: ErrInvalidValidation},
			},
		},
		{
			name: "invalid phone length",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345", "98765432109"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones[0]", Err: ErrInvalidValidation},
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid app version length",
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrInvalidValidation},
			},
		},
		{
			name: "valid token (no validation)",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrInvalidValidation},
			},
		},
		{
			name:        "not a struct",
			in:          "just a string",
			expectedErr: ErrNotStruct,
		},
		{
			name: "multiple validation errors",
			in: User{
				ID:     "short",
				Name:   "John Doe",
				Age:    15,
				Email:  "invalid",
				Role:   "guest",
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrInvalidValidation},
				{Field: "Age", Err: ErrInvalidValidation},
				{Field: "Email", Err: ErrInvalidValidation},
				{Field: "Role", Err: ErrInvalidValidation},
				{Field: "Phones[0]", Err: ErrInvalidValidation},
				{Field: "Phones[1]", Err: ErrInvalidValidation},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedErr)
				return
			}

			var expectedValidationErrs ValidationErrors
			if errors.As(tt.expectedErr, &expectedValidationErrs) {
				var actualValidationErrs ValidationErrors
				if !errors.As(err, &actualValidationErrs) {
					t.Errorf("expected validation errors, got %T", err)
					return
				}

				if len(actualValidationErrs) != len(expectedValidationErrs) {
					t.Errorf("expected %d validation errors, got %d", len(expectedValidationErrs), len(actualValidationErrs))
					return
				}

				for i := range expectedValidationErrs {
					if actualValidationErrs[i].Field != expectedValidationErrs[i].Field {
						t.Errorf("error %d: expected field %s, got %s", i, expectedValidationErrs[i].Field, actualValidationErrs[i].Field)
					}
					if !errors.Is(actualValidationErrs[i].Err, expectedValidationErrs[i].Err) {
						t.Errorf("error %d: expected error %v, got %v", i, expectedValidationErrs[i].Err, actualValidationErrs[i].Err)
					}
				}
			} else {
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			}
		})
	}
}
