package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError представляет ошибку валидации для конкретного поля.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors представляет список ошибок валидации.
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("%s: %v; ", err.Field, err.Err))
	}
	return sb.String()
}

var (
	ErrNotStruct         = errors.New("value is not a struct")
	ErrInvalidValidator  = errors.New("invalid validator tag")
	ErrInvalidRegexp     = errors.New("invalid regexp pattern")
	ErrInvalidValidation = errors.New("invalid validation rule")
)

// Validate проверяет структуру согласно тегам validate.
func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errs ValidationErrors
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		if !field.IsExported() {
			continue
		}

		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		if err := validateField(field, fieldValue, tag); err != nil {
			if validationErrs, ok := err.(ValidationErrors); ok { //nolint:errorlint
				errs = append(errs, validationErrs...)
			} else {
				return err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(field reflect.StructField, value reflect.Value, tag string) error {
	if tag == "nested" {
		return validateNestedField(field, value)
	}

	if value.Kind() == reflect.Slice {
		return validateSliceField(field, value, tag)
	}

	return validateSingleField(field.Name, value, tag)
}

func validateNestedField(field reflect.StructField, value reflect.Value) error {
	if value.Kind() != reflect.Struct {
		return nil
	}

	if nestedErr := Validate(value.Interface()); nestedErr != nil {
		var nestedValidationErrs ValidationErrors
		if errors.As(nestedErr, &nestedValidationErrs) {
			for _, nestedErr := range nestedValidationErrs {
				nestedErr.Field = field.Name + "." + nestedErr.Field
			}
			return nestedValidationErrs
		}
		return nestedErr
	}
	return nil
}

func validateSliceField(field reflect.StructField, value reflect.Value, tag string) error {
	var errs ValidationErrors
	for j := 0; j < value.Len(); j++ {
		element := value.Index(j)
		if err := validateSingleField(field.Name, element, tag); err != nil {
			if validationErrs, ok := err.(ValidationErrors); ok { //nolint:errorlint
				for _, validationErr := range validationErrs {
					validationErr.Field = fmt.Sprintf("%s[%d]", field.Name, j)
					errs = append(errs, validationErr)
				}
			} else {
				return err
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateSingleField(fieldName string, value reflect.Value, tag string) error {
	var errs ValidationErrors
	rules := strings.Split(tag, "|")

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return ErrInvalidValidator
		}

		validatorName := parts[0]
		validatorValue := parts[1]

		var err error
		switch value.Kind() { //nolint:exhaustive
		case reflect.String:
			err = validateString(value.String(), validatorName, validatorValue)
		case reflect.Int:
			err = validateInt(int(value.Int()), validatorName, validatorValue)
		default:
			return ErrInvalidValidator
		}

		if err != nil {
			if !errors.Is(err, ErrInvalidValidation) {
				return err
			}
			errs = append(errs, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateString(value, validatorName, validatorValue string) error {
	switch validatorName {
	case "len":
		expectedLength, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if len(value) != expectedLength {
			return ErrInvalidValidation
		}
	case "regexp":
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return ErrInvalidRegexp
		}
		if !re.MatchString(value) {
			return ErrInvalidValidation
		}
	case "in":
		options := strings.Split(validatorValue, ",")
		found := false
		for _, opt := range options {
			if value == opt {
				found = true
				break
			}
		}
		if !found {
			return ErrInvalidValidation
		}
	default:
		return ErrInvalidValidator
	}
	return nil
}

func validateInt(value int, validatorName, validatorValue string) error {
	switch validatorName {
	case "min":
		minValue, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if value < minValue {
			return ErrInvalidValidation
		}
	case "max":
		maxValue, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if value > maxValue {
			return ErrInvalidValidation
		}
	case "in":
		options := strings.Split(validatorValue, ",")
		found := false
		for _, opt := range options {
			num, err := strconv.Atoi(opt)
			if err != nil {
				return ErrInvalidValidator
			}
			if value == num {
				found = true
				break
			}
		}
		if !found {
			return ErrInvalidValidation
		}
	default:
		return ErrInvalidValidator
	}
	return nil
}
