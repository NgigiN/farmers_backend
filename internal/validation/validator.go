package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("safe_name", validateSafeName)
	validate.RegisterValidation("safe_location", validateSafeLocation)
}

func validateSafeName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsSafeName(value)
}

func validateSafeLocation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}
	return IsSafeLocation(value)
}

func ValidateStruct(obj any) error {
	return validate.Struct(obj)
}

func FormatError(err error) string {
	if err == nil {
		return ""
	}

	var validationErrs validator.ValidationErrors
	if !AsValidationErrors(err, &validationErrs) {
		return err.Error()
	}

	messages := make([]string, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		messages = append(messages, formatFieldError(fieldErr))
	}
	return strings.Join(messages, "; ")
}

func AsValidationErrors(err error, target *validator.ValidationErrors) bool {
	if err == nil {
		return false
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		*target = errs
		return true
	}
	return false
}

func formatFieldError(fieldErr validator.FieldError) string {
	field := strings.ToLower(fieldErr.Field())
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "email must be a valid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fieldErr.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fieldErr.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldErr.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fieldErr.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fieldErr.Param())
	case "safe_name":
		return fmt.Sprintf("%s contains invalid characters", field)
	case "safe_location":
		return fmt.Sprintf("%s contains invalid characters", field)
	default:
		return fmt.Sprintf("%s failed validation (%s)", field, fieldErr.Tag())
	}
}

func IsValidationError(err error) bool {
	var validationErrs validator.ValidationErrors
	return AsValidationErrors(err, &validationErrs)
}

func IsClientError(err error) bool {
	if err == nil {
		return false
	}
	if IsValidationError(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	clientMarkers := []string{
		"required",
		"not found",
		"does not belong",
		"invalid",
		"must be",
		"cannot",
		"password",
	}
	for _, marker := range clientMarkers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}

func DerefString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func StringPtr(value string) *string {
	if value == "" {
		return nil
	}
	v := value
	return &v
}

func sanitizeStringField(v reflect.Value) {
	if v.Kind() != reflect.String {
		return
	}
	v.SetString(SanitizeText(v.String()))
}