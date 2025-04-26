package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			assertValidationResult(t, tt.expectedErr, err)
		})
	}
}

func assertValidationResult(t *testing.T, expectedErr, actualErr error) {
	t.Helper()

	if expectedErr == nil {
		if actualErr != nil {
			t.Errorf("expected no error, got %v", actualErr)
		}
		return
	}

	if actualErr == nil {
		t.Errorf("expected error %v, got nil", expectedErr)
		return
	}

	if !assertValidationErrors(t, expectedErr, actualErr) {
		assertSimpleError(t, expectedErr, actualErr)
	}
}

func assertValidationErrors(t *testing.T, expectedErr, actualErr error) bool {
	t.Helper()

	var expectedValidationErrs ValidationErrors
	if !errors.As(expectedErr, &expectedValidationErrs) {
		return false
	}

	var actualValidationErrs ValidationErrors
	if !errors.As(actualErr, &actualValidationErrs) {
		t.Errorf("expected validation errors, got %T", actualErr)
		return true
	}

	if len(actualValidationErrs) != len(expectedValidationErrs) {
		t.Errorf("expected %d validation errors, got %d",
			len(expectedValidationErrs), len(actualValidationErrs))
		return true
	}

	for i := range expectedValidationErrs {
		if actualValidationErrs[i].Field != expectedValidationErrs[i].Field {
			t.Errorf("error %d: expected field %s, got %s",
				i, expectedValidationErrs[i].Field, actualValidationErrs[i].Field)
		}
		if !errors.Is(actualValidationErrs[i].Err, expectedValidationErrs[i].Err) {
			t.Errorf("error %d: expected error %v, got %v",
				i, expectedValidationErrs[i].Err, actualValidationErrs[i].Err)
		}
	}

	return true
}

func assertSimpleError(t *testing.T, expectedErr, actualErr error) {
	t.Helper()

	if !errors.Is(actualErr, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, actualErr)
	}
}
