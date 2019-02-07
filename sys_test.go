package bigip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SysTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *SysTestSuite) SetupSuite() {
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

func (s *SysTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *SysTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func (s *SysTestSuite) requireReserializesTo(expected string, actual interface{}, message string) {
	b, err := json.Marshal(actual)
	s.Require().Nil(err, message)

	s.Require().JSONEq(expected, string(b), message)
}

func TestSysSuite(t *testing.T) {
	suite.Run(t, new(SysTestSuite))
}

func (s *SysTestSuite) TestFolders() {
	resp := `{
  "items": [
    {
      "name": "/",
      "fullPath": "/",
      "deviceGroup": "none",
      "hidden": "false",
      "inheritedDevicegroup": "false",
      "inheritedTrafficGroup": "false",
      "noRefCheck": "false",
      "trafficGroup": "/Common/traffic-group-1"
    },
    {
      "name": "Common",
      "subPath": "/",
      "fullPath": "/Common",
      "deviceGroup": "none",
      "hidden": "false",
      "inheritedDevicegroup": "true",
      "inheritedTrafficGroup": "true",
      "noRefCheck": "false",
      "trafficGroup": "/Common/traffic-group-1"
    },
    {
      "name": "Drafts",
      "partition": "Common",
      "fullPath": "/Common/Drafts",
      "deviceGroup": "none",
      "hidden": "false",
      "inheritedDevicegroup": "true",
      "inheritedTrafficGroup": "true",
      "noRefCheck": "false",
      "trafficGroup": "/Common/traffic-group-1"
    },
    {
      "name": "test",
      "partition": "Common",
      "fullPath": "/Common/test",
      "deviceGroup": "none",
      "hidden": "false",
      "inheritedDevicegroup": "true",
      "inheritedTrafficGroup": "true",
      "noRefCheck": "false",
      "trafficGroup": "/Common/traffic-group-1"
    }
  ]
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	folders, err := s.Client.Folders()

	require.Nil(s.T(), err, "Error loading folders")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriSys, uriFolder), s.LastRequest.URL.Path, "Wrong uri to fetch folders")
	assert.Equal(s.T(), 4, len(folders.Folders), "Wrong number of folders")
	assert.Equal(s.T(), "/", folders.Folders[0].Name)
	assert.Equal(s.T(), "Common", folders.Folders[1].Name)
	assert.Equal(s.T(), "Drafts", folders.Folders[2].Name)
	assert.Equal(s.T(), "test", folders.Folders[3].Name)

	s.requireReserializesTo(resp, folders, "Folders should reserialize to itself")
}

func (s *SysTestSuite) TestCreateFolder() {
	err := s.Client.CreateFolder("/Common/test")

	require.Nil(s.T(), err, "Error creating folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriSys, uriFolder), s.LastRequest.URL.Path, "Wrong uri to create folders")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"/Common/test"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestAddFolder() {
	folder := Folder{
		Name:         "/Common/test",
		TrafficGroup: "default",
		NoRefCheck:   Bool(true),
	}
	err := s.Client.AddFolder(&folder)

	require.Nil(s.T(), err, "Error adding folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriSys, uriFolder), s.LastRequest.URL.Path, "Wrong uri to create folders")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"/Common/test", "trafficGroup": "default", "noRefCheck": "true"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestGetFolder() {
	resp := `{
  "name": "test",
  "partition": "Common",
  "fullPath": "/Common/test",
  "deviceGroup": "none",
  "hidden": "false",
  "inheritedDevicegroup": "true",
  "inheritedTrafficGroup": "true",
  "noRefCheck": "false",
  "trafficGroup": "/Common/traffic-group-1"
  }`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	folder, err := s.Client.GetFolder("/Common/test")

	require.Nil(s.T(), err, "Error getting folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFolder, "~Common~test"), s.LastRequest.URL.Path, "Wrong uri to fetch folders")
	s.requireReserializesTo(resp, folder, "Folder should reserialize to itself")
}

func (s *SysTestSuite) TestDeleteFolder() {
	err := s.Client.DeleteFolder("/Common/test")

	require.Nil(s.T(), err, "Error deleting folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFolder, "~Common~test"), s.LastRequest.URL.Path, "Wrong uri to create folders")
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *SysTestSuite) TestModifyFolder() {
	folder := Folder{
		TrafficGroup: "default",
		NoRefCheck:   Bool(true),
	}
	err := s.Client.ModifyFolder("/Common/test", &folder)

	require.Nil(s.T(), err, "Error modifying folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFolder, "~Common~test"), s.LastRequest.URL.Path, "Wrong uri to fetch folders")
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
}

func (s *SysTestSuite) TestPatchFolder() {
	folder := Folder{
		TrafficGroup: "default",
		NoRefCheck:   Bool(true),
	}
	err := s.Client.PatchFolder("/Common/test", &folder)

	require.Nil(s.T(), err, "Error patching folder")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFolder, "~Common~test"), s.LastRequest.URL.Path, "Wrong uri to fetch folders")
	assert.Equal(s.T(), "PATCH", s.LastRequest.Method)
}

func (s *SysTestSuite) TestCertificates() {
	resp := `{
	"items": [
        {
            "name": "/Common/foo.crt",
            "fullPath": "/Common/foo.crt",
            "generation": 1,
            "apiRawValues": {
                "certificateKeySize": "2048",
                "expiration": "Jan 01 00:00:00 2050 GMT",
                "publicKeyType": "RSA"
            },
            "city": "New York",
            "commonName": "foo.example.com",
            "country": "US",
            "emailAddress": "root@foo.example.com",
            "organization": "Foo Inc.",
            "ou": "IT",
            "state": "NY"
        },
		{
            "name": "/Common/bar.crt",
            "fullPath": "/Common/bar.crt",
            "generation": 1,
            "apiRawValues": {
                "certificateKeySize": "2048",
                "expiration": "Jan 01 00:00:00 2050 GMT",
                "publicKeyType": "RSA"
            },
            "city": "Los Angeles",
            "commonName": "bar.example.com",
            "country": "US",
            "emailAddress": "root@bar.example.com",
            "organization": "Bar Inc.",
            "ou": "IT",
            "state": "CA"
        }
	]
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	certs, err := s.Client.Certificates()

	require.Nil(s.T(), err, "Error loading certificates")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriCrypto, uriCert), s.LastRequest.URL.Path, "Wrong uri to fetch certificates")
	assert.Equal(s.T(), 2, len(certs.Certificates), "Wrong number of certificates")
	assert.Equal(s.T(), "/Common/foo.crt", certs.Certificates[0].Name)
	assert.Equal(s.T(), "/Common/bar.crt", certs.Certificates[1].Name)
	s.requireReserializesTo(resp, certs, "Certificates should reserialize to itself")
}

func (s *SysTestSuite) TestAddCertificate() {
	cert := Certificate{
		Name:          "test",
		Command:       "install",
		FromLocalFile: "/var/config/rest/downloads/test.crt",
	}
	err := s.Client.AddCertificate(&cert)

	require.Nil(s.T(), err, "Error adding certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriCrypto, uriCert), s.LastRequest.URL.Path, "Wrong uri to create certificate")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"test","command":"install","from-local-file":"/var/config/rest/downloads/test.crt"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestGetCertificate() {
	resp := `{
        "name": "/Common/test.crt",
        "fullPath": "/Common/test.crt",
        "generation": 1,
        "apiRawValues": {
            "certificateKeySize": "2048",
            "expiration": "Jan 01 00:00:00 2050 GMT",
            "publicKeyType": "RSA"
        },
        "city": "New York",
        "commonName": "test.example.com",
        "country": "US",
        "emailAddress": "root@test.example.com",
        "organization": "Test Inc.",
        "ou": "IT",
        "state": "NY"
	}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	cert, err := s.Client.GetCertificate("/Common/test.crt")

	require.Nil(s.T(), err, "Error getting certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriCrypto, uriCert, "~Common~test.crt"), s.LastRequest.URL.Path, "Wrong uri to fetch certificate")
	s.requireReserializesTo(resp, cert, "Certificate should reserialize to itself")
}

func (s *SysTestSuite) TestDeleteCertificate() {
	err := s.Client.DeleteCertificate("/Common/test.crt")

	require.Nil(s.T(), err, "Error deleting certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriCrypto, uriCert, "~Common~test.crt"), s.LastRequest.URL.Path, "Wrong uri to delete certificate")
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *SysTestSuite) TestKeys() {
	resp := `{
	"items": [
        {
            "name": "/Common/foo.key",
            "fullPath": "/Common/foo.key",
            "generation": 1,
            "keySize": "2048",
            "keyType": "rsa-private",
            "securityType": "normal"
        },
		{
            "name": "/Common/bar.key",
            "fullPath": "/Common/bar.key",
            "generation": 1,
            "keySize": "2048",
            "keyType": "rsa-private",
            "securityType": "normal"
        }
	]
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	keys, err := s.Client.Keys()

	require.Nil(s.T(), err, "Error loading keys")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriCrypto, uriKey), s.LastRequest.URL.Path, "Wrong uri to fetch keys")
	assert.Equal(s.T(), 2, len(keys.Keys), "Wrong number of keys")
	assert.Equal(s.T(), "/Common/foo.key", keys.Keys[0].Name)
	assert.Equal(s.T(), "/Common/bar.key", keys.Keys[1].Name)
	s.requireReserializesTo(resp, keys, "Keys should reserialize to itself")
}

func (s *SysTestSuite) TestAddKey() {
	key := Key{
		Name:          "test",
		Command:       "install",
		FromLocalFile: "/var/config/rest/downloads/test.key",
	}
	err := s.Client.AddKey(&key)

	require.Nil(s.T(), err, "Error adding key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriCrypto, uriKey), s.LastRequest.URL.Path, "Wrong uri to create key")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"test","command":"install","from-local-file":"/var/config/rest/downloads/test.key"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestGetKey() {
	resp := `{
		"name": "/Common/test.key",
		"fullPath": "/Common/test.key",
		"generation": 1,
		"keySize": "2048",
		"keyType": "rsa-private",
		"securityType": "normal"
	}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	key, err := s.Client.GetKey("/Common/test.key")

	require.Nil(s.T(), err, "Error getting key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriCrypto, uriKey, "~Common~test.key"), s.LastRequest.URL.Path, "Wrong uri to fetch key")
	s.requireReserializesTo(resp, key, "Key should reserialize to itself")
}

func (s *SysTestSuite) TestDeleteKey() {
	err := s.Client.DeleteKey("/Common/test.key")

	require.Nil(s.T(), err, "Error deleting key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriCrypto, uriKey, "~Common~test.key"), s.LastRequest.URL.Path, "Wrong uri to delete key")
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}
