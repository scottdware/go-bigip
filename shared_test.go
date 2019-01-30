package bigip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const validLicenseState = `{
  "vendor": "F5 Networks, Inc.",
  "licensedDateTime": "2018-04-02T00:00:00-07:00",
  "licensedVersion": "13.1.0",
  "licenseEndDateTime": "2018-05-03T00:00:00-07:00",
  "licenseStartDateTime": "2018-04-01T00:00:00-07:00",
  "registrationKey": "reg key here",
  "dossier": "dossier here",
  "authorization": "authorization here",
  "usage": "Evaluation",
  "platformId": "Z100k",
  "authVers": "5b",
  "serviceCheckDateTime": "2018-03-29T00:00:00-07:00",
  "machineId": "uuid here",
  "exclusivePlatform": [
    "Z100",
    "Z100A",
    "Z100a_icm",
    "Z100AzureCloud",
    "Z100GoogleCloud",
    "Z100K",
    "Z100x"
  ],
  "activeModules": [
    "Local Traffic Manager, VE-10G|RAYPVKR-EYOXDFQ|Routing Bundle, VE|Rate Shaping|APM, Limited|SSL, VE|Max Compression, VE|Anti-Virus Checks|Base Endpoint Security Checks|Firewall Checks|Network Access|Secure Virtual Keyboard|APM, Web Application|Machine Certificate Checks|Protected Workspace|Remote Desktop|App Tunnel"
  ],
  "optionalModules": [
    "Advanced Protocols, VE",
    "APC-VE, Introductory",
    "APC-VE, Introductory to Medium Upgrade",
    "APC-VE, Medium",
    "APM, Base, VE (50 CCU / 200 AS)",
    "App Mode (TMSH Only, No Root/Bash)",
    "BIG-IP VE, Multicast Routing",
    "CGN, ADD-VE, 10G",
    "DataSafe, VE-10G",
    "DNS and GTM (1K QPS), VE",
    "DNS and GTM (250 QPS), VE",
    "External Interface and Network HSM, VE",
    "FIPS 140-2 Level 1, BIG-IP VE-1G to 10G",
    "IP Intelligence, 1Yr, VE-10G",
    "IP Intelligence, 3Yr, VE-10G",
    "LTM to Best Bundle Upgrade, 10Gbps",
    "LTM to Better Bundle Upgrade, 10Gbps",
    "PEM, ADD-VE, 5G",
    "Secure Web Gateway, 1Yr, VE",
    "Secure Web Gateway, 3Yr, VE",
    "Secure Web Gateway, VE-3G-10G, 10000 Sessions, 1Yr",
    "Secure Web Gateway, VE-3G-10G, 10000 Sessions, 3Yr",
    "Secure Web Gateway, VE-3G-10G, 5000 Sessions, 1Yr",
    "Secure Web Gateway, VE-3G-10G, 5000 Sessions, 3Yr",
    "SSL, Forward Proxy, VE",
    "URL Filtering, 1Yr, VE",
    "URL Filtering, 3Yr, VE",
    "URL Filtering, VE-3G-10G, 10000 Sessions, 1Yr",
    "URL Filtering, VE-3G-10G, 10000 Sessions, 3Yr",
    "URL Filtering, VE-3G-10G, 5000 Sessions, 1Yr",
    "URL Filtering, VE-3G-10G, 5000 Sessions, 3Yr"
  ],
  "featureFlags": [
    {
      "featureName": "perf_SSL_Mbps",
      "featureValue": "1"
    },
    {
      "featureName": "apm_urlf_limited_sessions",
      "featureValue": "10"
    },
    {
      "featureName": "apml_sessions",
      "featureValue": "10"
    },
    {
      "featureName": "perf_VE_throughput_Mbps",
      "featureValue": "10000"
    },
    {
      "featureName": "perf_VE_cores",
      "featureValue": "8"
    },
    {
      "featureName": "perf_remote_crypto_client",
      "featureValue": "enabled"
    },
    {
      "featureName": "mod_ltm",
      "featureValue": "enabled"
    },
    {
      "featureName": "mod_ilx",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_network_virtualization",
      "featureValue": "enabled"
    },
    {
      "featureName": "perf_SSL_total_TPS",
      "featureValue": "UNLIMITED"
    },
    {
      "featureName": "perf_SSL_per_core",
      "featureValue": "enabled"
    },
    {
      "featureName": "perf_SSL_cmp",
      "featureValue": "enabled"
    },
    {
      "featureName": "perf_http_compression_Mbps",
      "featureValue": "UNLIMITED"
    },
    {
      "featureName": "nw_routing_rip",
      "featureValue": "enabled"
    },
    {
      "featureName": "nw_routing_ospf",
      "featureValue": "enabled"
    },
    {
      "featureName": "nw_routing_isis",
      "featureValue": "enabled"
    },
    {
      "featureName": "nw_routing_bgp",
      "featureValue": "enabled"
    },
    {
      "featureName": "nw_routing_bfd",
      "featureValue": "enabled"
    },
    {
      "featureName": "mod_apml",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_bandw_rate_tosque",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_bandw_rate_fairque",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_bandw_rate_classl7",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_bandw_rate_classl4",
      "featureValue": "enabled"
    },
    {
      "featureName": "ltm_bandw_rate_classes",
      "featureValue": "enabled"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "10.*.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.0.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.1.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.2.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.3.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.4.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "11.5.*"
    },
    {
      "featureName": "Deny_version",
      "featureValue": "9.*.*"
    },
    {
      "featureName": "apm_web_applications",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_remote_desktop",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_na",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep_svk",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep_pws",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep_machinecert",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep_fwcheck",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep_avcheck",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_ep",
      "featureValue": "enabled"
    },
    {
      "featureName": "apm_app_tunnel",
      "featureValue": "enabled"
    },
    {
      "featureName": "gtm_lc",
      "featureValue": "disabled"
    }
  ],
  "expiresInDays": "30.5",
  "expiresInDaysMessage": "License expires in 30 days, 11 hours."
}`

type SharedTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *SharedTestSuite) requireReserializesTo(expected string, actual interface{}, message string) {
	b, err := json.Marshal(actual)
	s.Require().Nil(err, message)

	s.Require().JSONEq(expected, string(b), message)
}

func (s *SharedTestSuite) SetupSuite() {
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		s.LastRequestBody = string(body)
		s.LastRequest = r
		if s.ResponseFunc != nil {
			s.ResponseFunc(w, r)
		}
	}))

	s.Client = NewSession(s.Server.URL, "", "", nil)
}

func (s *SharedTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *SharedTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestSharedSuite(t *testing.T) {
	suite.Run(t, new(SharedTestSuite))
}

func (s *SharedTestSuite) TestGetLicenseState() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(validLicenseState))
	}

	license, err := s.Client.GetLicenseState()

	s.Require().Nil(err, "Error getting license state")
	s.Require().NotNil(license, "License should not be nil")

	s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), s.LastRequest.URL.Path, "Wrong uri to fetch license")
	s.Require().Equal("reg key here", license.RegistrationKey)
	s.requireReserializesTo(validLicenseState, license, "License should reserialize to itself")
}

func (s *SharedTestSuite) TestInstallLicense() {
	err := s.Client.InstallLicense("this is a new license\nline 2")

	s.Require().Nil(err, "Error installing license")
	s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), s.LastRequest.URL.Path, "Wrong uri to install license")
	s.Require().Equal("PUT", s.LastRequest.Method, "Wrong method to install license")
	s.Require().JSONEq(`{"licenseText": "this is a new license\nline 2"}`, s.LastRequestBody)
}

func (s *SharedTestSuite) TestActivateTimeout() {
	err := s.Client.AutoLicense("reg key", nil, 0)
	s.Require().Error(err)
	s.Require().Equal("Timed out after 0s", err.Error(), "Should timeout")
}

func (s *SharedTestSuite) TestActivateNoEula() {
	counter := 0
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		counter += 1
		switch counter {
		case 1:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			s.Require().Equal("POST", r.Method, "Wrong method to activate license")
			w.Write([]byte(`{}`))
		case 2:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			w.Write([]byte(`{"status": "LICENSING_COMPLETE", "licenseText": "this is a new license\nline 2"}`))
		case 3:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), r.URL.Path, "Wrong uri to install license")
			w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Ran past number of expected requests"}`))
		}
	}

	err := s.Client.AutoLicense("reg key", nil, 10*time.Second)
	s.Require().Nil(err, "Error auto activating license")
	s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), s.LastRequest.URL.Path, "Wrong uri to install license")
	s.Require().Equal("PUT", s.LastRequest.Method, "Wrong method to install license")
	s.Require().JSONEq(`{"licenseText": "this is a new license\nline 2"}`, s.LastRequestBody)
}

func (s *SharedTestSuite) TestActivateEula() {
	counter := 0
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		counter += 1
		switch counter {
		case 1:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			s.Require().Equal("POST", r.Method, "Wrong method to activate license")
			w.Write([]byte(`{}`))
		case 2:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			w.Write([]byte(`{"status": "NEED_EULA_ACCEPT", "eulaText": "eula to accept"}`))
		case 3:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			s.Require().Equal("POST", r.Method, "Wrong method to install license")
			s.Require().JSONEq(`{"baseRegKey": "reg key", "eulaText": "eula to accept", "isAutomaticActivation": true}`, s.LastRequestBody)
			w.Write([]byte(`{}`))
		case 4:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriActivation), r.URL.Path, "Wrong uri to activate")
			w.Write([]byte(`{"status": "LICENSING_COMPLETE", "licenseText": "this is a new license\nline 2"}`))
		case 5:
			s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), r.URL.Path, "Wrong uri to install license")
			w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Ran past number of expected requests"}`))
		}
	}

	err := s.Client.AutoLicense("reg key", nil, 10*time.Second)
	s.Require().Nil(err, "Error auto activating license")
	s.Require().Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriShared, uriLicensing, uriRegistration), s.LastRequest.URL.Path, "Wrong uri to install license")
	s.Require().Equal("PUT", s.LastRequest.Method, "Wrong method to install license")
	s.Require().JSONEq(`{"licenseText": "this is a new license\nline 2"}`, s.LastRequestBody)
}

func (s *SharedTestSuite) TestUploadFile() {
	tmp, err := ioutil.TempFile("", "test")
	defer os.Remove(tmp.Name())
	content := []byte("test file content")
	size := len(content)
	tmp.Write(content)
	tmp.Close()
	f, _ := os.Open(tmp.Name())
	filename := path.Base(tmp.Name())

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	}

	upload, err := s.Client.UploadFile(f)
	s.Require().Nil(err, "Error uploading file")
	s.Require().NotNil(upload, "Upload response should not be nil")
	s.Require().Equal("application/octet-stream", s.LastRequest.Header.Get("Content-Type"), "Wrong Content-Type header")
	s.Require().Equal(fmt.Sprintf("0-%d/%d", size-1, size), s.LastRequest.Header.Get("Content-Range"), "Wrong Content-Range header")
	s.Require().Equal(fmt.Sprintf("/mgmt/shared/file-transfer/uploads/%s", filename), s.LastRequest.URL.Path, "Wrong uri to upload file")
}

func (s *SharedTestSuite) TestUploadBytes() {
	b := []byte("test byte content")
	size := len(b)
	filename := "test.txt"

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	}

	upload, err := s.Client.UploadBytes(b, filename)
	s.Require().Nil(err, "Error uploading file")
	s.Require().NotNil(upload, "Upload response should not be nil")
	s.Require().Equal("application/octet-stream", s.LastRequest.Header.Get("Content-Type"), "Wrong Content-Type header")
	s.Require().Equal(fmt.Sprintf("0-%d/%d", size-1, size), s.LastRequest.Header.Get("Content-Range"), "Wrong Content-Range header")
	s.Require().Equal(fmt.Sprintf("/mgmt/shared/file-transfer/uploads/%s", filename), s.LastRequest.URL.Path, "Wrong uri to upload file")
}
