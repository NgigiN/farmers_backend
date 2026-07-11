package validation

import (
	"errors"
	"strings"
)

var (
	ErrInvalidSourceType = errors.New("source_type must be either 'plant' or 'animal'")
	ErrInvalidSource     = errors.New("source must be either 'plant' or 'animal'")
	ErrInvalidCategory   = errors.New("category must be one of: input, activity, infrastructure")
	ErrInvalidCostType   = errors.New("type must be either 'plant' or 'animal'")
)

func NormalizeOptionalEnum(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func ValidateSourceType(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	normalized := NormalizeOptionalEnum(value)
	if normalized != "plant" && normalized != "animal" {
		return "", ErrInvalidSourceType
	}
	return normalized, nil
}

func ValidateRevenueSource(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	normalized := NormalizeOptionalEnum(value)
	if normalized != "plant" && normalized != "animal" {
		return "", ErrInvalidSource
	}
	return normalized, nil
}

func ValidateCostCategoryType(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	normalized := NormalizeOptionalEnum(value)
	if normalized != "plant" && normalized != "animal" {
		return "", ErrInvalidCostType
	}
	return normalized, nil
}

func ValidateCostCategoryCategory(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	normalized := NormalizeOptionalEnum(value)
	switch normalized {
	case "input", "activity", "infrastructure":
		return normalized, nil
	default:
		return "", ErrInvalidCategory
	}
}