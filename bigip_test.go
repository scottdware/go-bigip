package bigip

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestMarshalYesTrue(t *testing.T) {
	v := Yes(true)
	testMarshalBoolean(&v, "yes", t);
}
func TestMarshalYesFalse(t *testing.T) {
	v:=Yes(false)
	testMarshalBoolean(&v, "no", t);
}
func TestMarshalEnabledTrue(t *testing.T) {
	v:=Enabled(true)
	testMarshalBoolean(&v, "enabled", t);
}
func TestMarshalEnabledFalse(t *testing.T) {
	v:=Enabled(false)
	testMarshalBoolean(&v, "disabled", t);
}
func TestMarshalTrueTrue(t *testing.T) {
	v:=True(true)
	testMarshalBoolean(&v, "true", t);
}
func TestMarshalTrueFalse(t *testing.T) {
	v:=True(false)
	testMarshalBoolean(&v, "false", t);
}
func TestUnmarshalYesTrue(t *testing.T){
	var v Yes
	json.Unmarshal([]byte("\"yes\""), &v)
	if !v {
		t.Error("Expected true")
	}
}
func TestUnmarshalYesFalse(t *testing.T){
	var v Yes
	json.Unmarshal([]byte("\"no\""), &v)
	if v {
		t.Error("Expected false")
	}
}
func TestUnmarshalEnabledTrue(t *testing.T){
	var v Enabled
	json.Unmarshal([]byte("\"enabled\""), &v)
	if !v {
		t.Error("Expected true")
	}
}
func TestUnmarshalEnabledFalse(t *testing.T){
	var v Enabled
	json.Unmarshal([]byte("\"disabled\""), &v)
	if v {
		t.Error("Expected false")
	}
}
func TestUnmarshalTrueTrue(t *testing.T){
	var v True
	json.Unmarshal([]byte("\"true\""), &v)
	if !v {
		t.Error("Expected true")
	}
}
func TestUnmarshalTrueFalse(t *testing.T){
	var v True
	json.Unmarshal([]byte("\"false\""), &v)
	if v {
		t.Error("Expected false")
	}
}

func testMarshalBoolean(v json.Marshaler, expected string, t *testing.T) {
	b, _ := json.Marshal(v);
	json := fmt.Sprintf("%s", b)

	if json != "\"" + expected + "\"" {
		t.Error("Expected \"" + expected + "\" got: ", json)
	}
}