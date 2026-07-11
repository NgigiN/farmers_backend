package validation

import (
	"testing"
)

func TestSanitizeTextStripsControlCharacters(t *testing.T) {
	got := SanitizeText("  Farm\x00Name  ")
	if got != "FarmName" {
		t.Fatalf("expected FarmName, got %q", got)
	}
}

func TestRequiredNameRejectsInvalidCharacters(t *testing.T) {
	req := UpdateProfileRequest{
		FarmName: "<script>",
		Location: "Nairobi",
	}
	req.Sanitize()
	if err := ValidateStruct(&req); err == nil {
		t.Fatal("expected validation error for unsafe farm name")
	}
}

func TestLandRequestAcceptsValidInput(t *testing.T) {
	size := float32(2.5)
	req := LandRequest{
		Name:     "North Field",
		Size:     &size,
		Location: "Nairobi, Kenya",
		SoilType: "Loam",
	}
	req.Sanitize()
	if err := ValidateStruct(&req); err != nil {
		t.Fatalf("expected valid land request, got %v", err)
	}
}

func TestValidateSourceTypeRejectsInvalidValue(t *testing.T) {
	if _, err := ValidateSourceType("invalid"); err == nil {
		t.Fatal("expected invalid source_type error")
	}
}

func TestRegisterRequestRequiresPassword(t *testing.T) {
	req := RegisterRequest{
		Email:     "user@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
		Password:  "short",
	}
	if err := ValidateStruct(&req); err == nil {
		t.Fatal("expected password length validation error")
	}
}