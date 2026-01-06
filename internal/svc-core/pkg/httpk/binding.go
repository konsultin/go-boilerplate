package httpk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	f "github.com/valyala/fasthttp"
)

// Singleton validator instance for performance
var (
	validate     *validator.Validate
	validateOnce sync.Once
)

// getValidator returns a singleton validator instance
func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()

		// Use json tag names for error messages instead of struct field names
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return fld.Name
			}
			return name
		})

		// Register custom validations here if needed
		// Example: validate.RegisterValidation("custom_tag", customValidationFunc)
	})
	return validate
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	var messages []string
	for _, e := range v.Errors {
		messages = append(messages, e.Message)
	}
	return strings.Join(messages, "; ")
}

// formatValidationError creates a human-readable error message
func formatValidationError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number (E.164 format)", field)
	default:
		return fmt.Sprintf("%s failed on '%s' validation", field, fe.Tag())
	}
}

// Validate validates a struct using go-playground/validator
func Validate(payload any) error {
	v := getValidator()
	err := v.Struct(payload)
	if err == nil {
		return nil
	}

	// Type assert to validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var errors []ValidationError
	for _, fe := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   fe.Field(),
			Tag:     fe.Tag(),
			Value:   fe.Value(),
			Message: formatValidationError(fe),
		})
	}

	return ValidationErrors{Errors: errors}
}

// BindJSON binds JSON body to a struct
func BindJSON[T any](ctx *f.RequestCtx) (*T, error) {
	var payload T
	body := ctx.Request.Body()

	if len(body) == 0 {
		return nil, InvalidPayloadError.Wrap(fmt.Errorf("request body is empty"))
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, InvalidPayloadError.Wrap(err)
	}

	return &payload, nil
}

// BindAndValidate binds JSON body and validates the struct
// This is the recommended function for most use cases
func BindAndValidate[T any](ctx *f.RequestCtx) (*T, error) {
	payload, err := BindJSON[T](ctx)
	if err != nil {
		return nil, err
	}

	if err := Validate(payload); err != nil {
		return nil, InvalidPayloadError.Wrap(err)
	}

	return payload, nil
}

// BindQuery binds query parameters to a struct
func BindQuery[T any](ctx *f.RequestCtx) (*T, error) {
	var payload T
	v := reflect.ValueOf(&payload).Elem()
	t := v.Type()

	args := ctx.QueryArgs()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Get the query tag or fallback to json tag
		tag := field.Tag.Get("query")
		if tag == "" {
			tag = field.Tag.Get("json")
		}
		if tag == "" || tag == "-" {
			continue
		}

		// Remove omitempty and other options
		tagName := strings.SplitN(tag, ",", 2)[0]
		queryValue := string(args.Peek(tagName))

		if queryValue == "" {
			continue
		}

		// Set the value based on field type
		if err := setFieldValue(fieldValue, queryValue); err != nil {
			return nil, InvalidPayloadError.Wrap(fmt.Errorf("invalid value for field %s: %w", field.Name, err))
		}
	}

	return &payload, nil
}

// BindQueryAndValidate binds query parameters and validates the struct
func BindQueryAndValidate[T any](ctx *f.RequestCtx) (*T, error) {
	payload, err := BindQuery[T](ctx)
	if err != nil {
		return nil, err
	}

	if err := Validate(payload); err != nil {
		return nil, InvalidPayloadError.Wrap(err)
	}

	return payload, nil
}

// setFieldValue sets a reflect.Value from a string
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intVal int64
		if _, err := fmt.Sscanf(value, "%d", &intVal); err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var uintVal uint64
		if _, err := fmt.Sscanf(value, "%d", &uintVal); err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		var floatVal float64
		if _, err := fmt.Sscanf(value, "%f", &floatVal); err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		field.SetBool(value == "true" || value == "1")
	case reflect.Ptr:
		// Handle pointer types
		elem := reflect.New(field.Type().Elem())
		if err := setFieldValue(elem.Elem(), value); err != nil {
			return err
		}
		field.Set(elem)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
