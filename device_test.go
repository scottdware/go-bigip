package bigip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
)

var goodDeviceResponse = `{
    "items": [
        {
            "activeModules": [
                "APM, Base, VE GBB (500 CCU, 2500 Access Sessions)|XFLGBEL-BRZXIDZ|Anti-Virus Checks|Base Endpoint Security Checks|Firewall Checks|Network Access|Secure Virtual Keyboard|APM, Web Application|Machine Certificate Checks|Protected Workspace|Remote Desktop|App Tunnel"
            ],
            "baseMac": "aa:aa:aa:aa:aa:aa",
            "build": "0.0.4",
            "cert": "/Common/dtdi.crt",
            "certReference": {
                "link": "https://localhost/mgmt/tm/cm/cert/~Common~dtdi.crt?ver=13.1.1.2"
            },
            "chassisId": "foo",
            "chassisType": "individual",
            "configsyncIp": "10.1.1.10",
            "edition": "Point Release 2",
            "failoverState": "active",
            "fullPath": "/Common/foo.f5.com",
            "generation": 1,
            "haCapacity": 0,
            "hostname": "foo.f5.com",
            "key": "/Common/dtdi.key",
            "keyReference": {
                "link": "https://localhost/mgmt/tm/cm/key/~Common~dtdi.key?ver=13.1.1.2"
            },
            "kind": "tm:cm:device:devicestate",
            "managementIp": "172.16.1.10",
            "marketingName": "BIG-IP Virtual Edition",
            "mirrorIp": "10.1.1.10",
            "mirrorSecondaryIp": "any6",
            "multicastIp": "any6",
            "multicastPort": 0,
            "name": "foo.f5.com",
            "optionalModules": [
                "Advanced Protocols, VE",
                "URL Filtering, VE-25M-1G, 500 Sessions, 3Yr"
            ],
            "partition": "Common",
            "platformId": "foo",
            "product": "BIG-IP",
            "selfDevice": "true",
            "selfLink": "https://localhost/mgmt/tm/cm/device/~Common~foo.f5.com?ver=13.1.1.2",
            "timeZone": "Europe/Rome",
            "unicastAddress": [
                {
                    "effectiveIp": "management-ip",
                    "effectivePort": 123,
                    "ip": "management-ip",
                    "port": 123
                },
                {
                    "effectiveIp": "10.1.1.10",
                    "effectivePort": 456,
                    "ip": "10.1.1.10",
                    "port": 456
                }
            ],
            "version": "13.1.1.2"
        },
        {
            "activeModules": [
                "APM, Base, VE GBB (500 CCU, 2500 Access Sessions)|UEMOINR-IPCMIEL|Anti-Virus Checks|Base Endpoint Security Checks|Firewall Checks|Network Access|Secure Virtual Keyboard|APM, Web Application|Machine Certificate Checks|Protected Workspace|Remote Desktop|App Tunnel"
            ],
            "baseMac": "bb:bb:bb:bb:bb:bb",
            "build": "0.0.4",
            "chassisId": "bar",
            "chassisType": "individual",
            "configsyncIp": "10.1.1.11",
            "edition": "Point Release 2",
            "failoverState": "standby",
            "fullPath": "/Common/bar.f5.com",
            "generation": 2,
            "haCapacity": 0,
            "hostname": "bar.f5.com",
            "kind": "tm:cm:device:devicestate",
            "managementIp": "172.16.1.11",
            "marketingName": "BIG-IP Virtual Edition",
            "mirrorIp": "10.1.1.11",
            "mirrorSecondaryIp": "any6",
            "multicastIp": "any6",
            "multicastPort": 0,
            "name": "bar.f5.com",
            "optionalModules": [
                "Advanced Protocols, VE",
                "URL Filtering, VE-25M-1G, 500 Sessions, 3Yr"
            ],
            "partition": "Common",
            "platformId": "Z100",
            "product": "BIG-IP",
            "selfDevice": "false",
            "selfLink": "https://localhost/mgmt/tm/cm/device/~Common~bar.f5.com?ver=13.1.1.2",
            "timeZone": "Europe/Rome",
            "unicastAddress": [
                {
                    "effectiveIp": "management-ip",
                    "effectivePort": 123,
                    "ip": "management-ip",
                    "port": 123
                },
                {
                    "effectiveIp": "10.1.1.11",
                    "effectivePort": 456,
                    "ip": "10.1.1.11",
                    "port": 456
                }
            ],
            "version": "13.1.1.2"
        }
    ],
    "kind": "tm:cm:device:devicecollectionstate",
    "selfLink": "https://localhost/mgmt/tm/cm/device?ver=13.1.1.2"
}`

// assertDeviceRestCall verifies that the DeviceTestSuite receives an
// http request with a matching method, URI, and body.
func assertDeviceRestCall(s *DeviceTestSuite, method, path, body string) {
	assert.Equal(s.T(), method, s.LastRequest.Method)
	assert.Equal(s.T(), path, s.LastRequest.URL.Path)
	if body != "" {
		assert.JSONEq(s.T(), body, s.LastRequestBody)
	}
}

// Setup the TestSuite

type DeviceTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *DeviceTestSuite) SetupSuite() {
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		s.LastRequestBody = string(body)
		s.LastRequest = r
		if s.ResponseFunc != nil {
			s.ResponseFunc(w, r)
		}
	}))

	s.Client = NewSession(s.Server.URL, "", "", "", nil)
}

func (s *DeviceTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *DeviceTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceTestSuite))
}

func (s *DeviceTestSuite) TestGetDevices() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(goodDeviceResponse))
	}

	devices, err := s.Client.GetDevices()

	assert.Nil(s.T(), err)
	assertDeviceRestCall(s, "GET", "/mgmt/tm/cm/device", "")
	assert.Equal(s.T(), 2, len(devices))

	assert.Equal(s.T(), "foo.f5.com", devices[0].Name)
	assert.Equal(s.T(), "bar.f5.com", devices[1].Name)

	assert.Equal(s.T(), "active", devices[0].FailoverState)
	assert.Equal(s.T(), "standby", devices[1].FailoverState)
}
