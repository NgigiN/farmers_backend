package middleware

import "farm-backend/internal/validation"

func ValidateStruct(obj any) error {
	return validation.ValidateStruct(obj)
}