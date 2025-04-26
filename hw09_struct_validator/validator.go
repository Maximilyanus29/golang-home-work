package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError представляет ошибку валидации для конкретного поля
type ValidationError struct {
	Field string
	Err   error
}

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

		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		if tag == "nested" {
			if fieldValue.Kind() == reflect.Struct {
				if nestedErr := Validate(fieldValue.Interface()); nestedErr != nil {
					var nestedValidationErrs ValidationErrors
					if errors.As(nestedErr, &nestedValidationErrs) {
						for _, nestedErr := range nestedValidationErrs {
							nestedErr.Field = field.Name + "." + nestedErr.Field
							errs = append(errs, nestedErr)
						}
					} else {
						return nestedErr
					}
				}
			}
			continue
		}

		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < fieldValue.Len(); j++ {
				element := fieldValue.Index(j)
				if err := validateField(field.Name, element, tag); err != nil {
					if validationErrs, ok := err.(ValidationErrors); ok {
						for _, validationErr := range validationErrs {
							validationErr.Field = fmt.Sprintf("%s[%d]", field.Name, j)
							errs = append(errs, validationErr)
						}
					} else {
						return err
					}
				}
			}
		} else {
			if err := validateField(field.Name, fieldValue, tag); err != nil {
				if validationErrs, ok := err.(ValidationErrors); ok {
					errs = append(errs, validationErrs...)
				} else {
					return err
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(fieldName string, value reflect.Value, tag string) error {
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
		switch value.Kind() {
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
		expectedLen, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if len(value) != expectedLen {
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
		min, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if value < min {
			return ErrInvalidValidation
		}
	case "max":
		max, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidator
		}
		if value > max {
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
