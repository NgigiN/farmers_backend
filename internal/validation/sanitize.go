package validation

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	safeNamePattern     = regexp.MustCompile(`^[a-zA-Z0-9\s\-'.]+$`)
	safeLocationPattern = regexp.MustCompile(`^[a-zA-Z0-9\s\-'.,]+$`)
)

func Trim(value string) string {
	return strings.TrimSpace(value)
}

func StripControlChars(value string) string {
	if value == "" {
		return value
	}
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		if r == '\t' || r == '\n' || r == '\r' || !unicode.IsControl(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func SanitizeText(value string) string {
	return StripControlChars(Trim(value))
}

func SanitizeOptionalText(value string) string {
	sanitized := SanitizeText(value)
	if sanitized == "" {
		return ""
	}
	return sanitized
}

func IsSafeName(value string) bool {
	return safeNamePattern.MatchString(value)
}

func IsSafeLocation(value string) bool {
	return safeLocationPattern.MatchString(value)
}