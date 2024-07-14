package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	UnexpectedErr  = "Unexpected"
	ErrInvalidType = func(field string, expectedType any, err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidType",
			message:   fmt.Sprintf("InvalidType for field: %s. Expected: %s", field, expectedType),
			Err:       err,
		}
	}
	ErrInvalidJson = func(err error) ValidationErrorInterface {
		return &ValidationError{
			errorType: "InvalidJson",
			message:   fmt.Sprintf("InvalidJson: %s", err.Error()),
			Err:       err,
		}
	}
	ErrInvalidValue = func(message string, err error) ValidationErrorInterface {
		if message == "" {
			message = err.Error()
		}
		return &ValidationError{
			errorType: "InvalidValue",
			message:   fmt.Sprintf("InvalidValue: %s", message),
			Err:       err,
		}
	}
)

type ValidationErrorInterface interface {
	Type() string
	Error() string
	Unwrap() error
	IsUnexpectedErr() bool
}

type ValidationError struct {
	errorType string
	message   string
	Err       error
}

func (e *ValidationError) Error() string { return e.message }

func (e *ValidationError) Type() string { return e.errorType }

func (e *ValidationError) Unwrap() error { return e.Err }

func (e *ValidationError) IsUnexpectedErr() bool { return e.errorType == UnexpectedErr }

func HandleValidationErrors(err error) ValidationErrorInterface {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		return ErrInvalidType(e.Field, e.Type, e)
	case validator.ValidationErrors:
		var msgs []string
		for _, fe := range e {
			switch fe.Tag() {
			case "required":
				msgs = append(msgs, fmt.Sprintf("%s is a required field", fe.Field()))
			case "notblank":
				msgs = append(msgs, fmt.Sprintf("%s should not be empty", fe.Field()))
			case "min":
				msgs = append(msgs, fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param()))
			case "max":
				msgs = append(msgs, fmt.Sprintf("%s must be a maximum of %s in length", fe.Field(), fe.Param()))
			case "url":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", fe.Field()))
			case "uuid":
				msgs = append(msgs, fmt.Sprintf("%s must be a valid uuid", fe.Field()))
			case "oneof":
				msgs = append(msgs, fmt.Sprintf("%s value must be one of the pre-decided values", fe.Field()))
			case "isodateformat":
				msgs = append(msgs, fmt.Sprintf("%s must be in yyyy-mm-dd format", fe.Field()))
			case "phone_number":
				msgs = append(msgs, fmt.Sprintf("Phone number is invalid %s", fe.Value()))
			case "email":
				msgs = append(msgs, fmt.Sprintf("Email is invalid %s", fe.Value()))
			default:
				msgs = append(msgs, fmt.Sprintf("validation failed for %s on %s", fe.Field(), fe.Tag()))
			}
		}
		return ErrInvalidValue(strings.Join(msgs, ", "), e)
	case *json.SyntaxError:
		return ErrInvalidJson(e)
	default:
		return &ValidationError{
			errorType: UnexpectedErr,
			message:   e.Error(),
			Err:       e,
		}
	}
}

func ValidateStruct(s any, structValidations ...validator.StructLevelFunc) ValidationErrorInterface {
	var validate = validator.New()
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	return HandleValidationErrors(err)
}
