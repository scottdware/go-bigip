package f5teem

import (
	"testing"
)

func TestTeemTelemetryRequest(t *testing.T) {
	assetInfo := AssetInfo{
		"Terraform-Provider-BIGIP-Ecosystem",
		"1.2.0",
		"",
	}
	teemDevice := AnonymousClient(assetInfo, "")
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
	if err != nil {
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
