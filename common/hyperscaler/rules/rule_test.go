package rules

import (
	"fmt"
	"testing"
)

func TestSignatures(t *testing.T) {
	rule := NewRule()
	rule.SetPlanNoValidation("gcp")
	rule.PlatformRegionSuffix = true
	rule.PlatformRegion = "cf-eu30"

	fmt.Println(rule.SignatureWithValues())
	fmt.Println(rule.MirroredSignature())
}
