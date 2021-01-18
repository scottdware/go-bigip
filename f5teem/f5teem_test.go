package f5teem

import (
	"os"
	"testing"
)

func TestTeemTelemetryRequest(t *testing.T) {
	assetInfo := AssetInfo{
		"Terraform-Provider-BIGIP-Ecosystem",
		"1.2.0",
		"",
	}
	apiKey := os.Getenv("TEEM_API_KEY")
	teemDevice := AnonymousClient(assetInfo, apiKey)
	d := map[string]interface{}{
		"Device":          1,
		"Tenant":          1,
		"License":         1,
		"DNS":             1,
		"NTP":             1,
		"Provision":       1,
		"VLAN":            2,
		"SelfIp":          2,
		"platform":        "BIG-IP",
		"platformVersion": "15.1.0.5",
	}
	err := teemDevice.Report(d, "Terraform BIGIP-ravinder-latest", "1")
	if apiKey == "" && err == nil {
		t.Errorf("Error:%v", err)
	}
	if apiKey != "" && err != nil {
		t.Errorf("Error:%v", err)
	}
}

func TestTeemNotAuthorized(t *testing.T) {
	assetInfo := AssetInfo{
		"Terraform-Provider-BIGIP-Ecosystem",
		"1.2.0",
		"",
	}
	teemDevice := AnonymousClient(assetInfo, "xxxx")
	d := map[string]interface{}{
		"Device":          1,
		"Tenant":          1,
		"License":         1,
		"DNS":             1,
		"NTP":             1,
		"Provision":       1,
		"VLAN":            2,
		"SelfIp":          2,
		"platform":        "BIG-IP",
		"platformVersion": "15.1.0.5",
	}
	err := teemDevice.Report(d, "Terraform BIGIP-ravinder-latest", "1")
	if err == nil {
		t.Errorf("Error:%v", err)
	}
}

func TestUniqueUUID(t *testing.T) {
	expected := "a837ce9d-d34c-e5c1-9fe0-581e2b46c029"
	oldOS := osHostname
	defer func() { osHostname = oldOS }()

	osHostname = func() (hostname string, err error) {
		return "foobar.local.lab", nil
	}
	result := uniqueUUID()
	if result != expected {
		t.Errorf("Expected UUID to be: %s, got %s", expected, result)
	}
}
