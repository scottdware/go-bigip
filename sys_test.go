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
      "name": "foo.crt",
      "partition": "Common",
      "fullPath": "/Common/foo.crt",
      "generation": 1,
      "certificateKeyCurveName": "none",
      "certificateKeySize": 2048,
      "checksum": "SHA1:1184:d8af9b644095d0c0626ad38f1cec5b9e4188e794",
      "createTime": "2019-03-15T15:10:11Z",
      "createdBy": "admin",
      "expirationDate": 1584198611,
      "expirationString": "Mar 14 15:10:11 2020 GMT",
      "isBundle": "false",
      "issuer": "CN=foo.example.com,O=Foo Inc.,L=New York,ST=New York,C=US",
      "keyType": "rsa-public",
      "lastUpdateTime": "2019-03-15T15:10:11Z",
      "mode": 33188,
      "revision": 1,
      "serialNumber": "291519056",
      "size": 1184,
      "sourcePath": "/config/ssl/ssl.crt/foo.crt",
      "subject": "CN=foo.example.com,O=Foo Inc.,L=New York,ST=New York,C=US",
      "updatedBy": "admin",
      "version": 1
    },
    {
      "name": "bar.crt",
      "partition": "Common",
      "fullPath": "/Common/bar.crt",
      "generation": 1,
      "certificateKeyCurveName": "none",
      "certificateKeySize": 2048,
      "checksum": "SHA1:1196:7b55ca0b4dc10bdb29a4bca5d0d70627e3fc63f3",
      "createTime": "2019-03-15T15:10:31Z",
      "createdBy": "admin",
      "expirationDate": 1584198631,
      "expirationString": "Mar 14 15:10:31 2020 GMT",
      "isBundle": "false",
      "issuer": "CN=bar.example.com,O=Bar Inc.,L=Los Angeles,ST=California,C=US",
      "keyType": "rsa-public",
      "lastUpdateTime": "2019-03-15T15:10:31Z",
      "mode": 33188,
      "revision": 1,
      "serialNumber": "203165026",
      "size": 1196,
      "sourcePath": "/config/ssl/ssl.crt/bar.crt",
      "subject": "CN=bar.example.com,O=Bar Inc.,L=Los Angeles,ST=California,C=US",
      "updatedBy": "admin",
      "version": 1
    }
  ]
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	certs, err := s.Client.Certificates()

	require.Nil(s.T(), err, "Error loading certificates")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFile, uriSslCert), s.LastRequest.URL.Path, "Wrong uri to fetch certificates")
	assert.Equal(s.T(), 2, len(certs.Certificates), "Wrong number of certificates")
	assert.Equal(s.T(), "foo.crt", certs.Certificates[0].Name)
	assert.Equal(s.T(), "bar.crt", certs.Certificates[1].Name)
	s.requireReserializesTo(resp, certs, "Certificates should reserialize to itself")
}

func (s *SysTestSuite) TestAddCertificate() {
	cert := Certificate{
		Name:       "test",
		SourcePath: "file:///var/config/rest/downloads/test.crt",
	}
	err := s.Client.AddCertificate(&cert)

	require.Nil(s.T(), err, "Error adding certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFile, uriSslCert), s.LastRequest.URL.Path, "Wrong uri to create certificate")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"test", "sourcePath":"file:///var/config/rest/downloads/test.crt"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestGetCertificate() {
	resp := `{
  "name": "test.crt",
  "partition": "Common",
  "fullPath": "/Common/test.crt",
  "generation": 1,
  "certificateKeyCurveName": "none",
  "certificateKeySize": 2048,
  "checksum": "SHA1:1188:ef0223d316fb0f16e07fc25e4c3f396ff4f43e8e",
  "createTime": "2019-03-15T15:09:46Z",
  "createdBy": "admin",
  "expirationDate": 1584198586,
  "expirationString": "Mar 14 15:09:46 2020 GMT",
  "isBundle": "false",
  "issuer": "CN=test.example.com,O=Test Inc.,L=New York,ST=New York,C=US",
  "keyType": "rsa-public",
  "lastUpdateTime": "2019-03-15T15:09:46Z",
  "mode": 33188,
  "revision": 1,
  "serialNumber": "219156055",
  "size": 1188,
  "sourcePath": "/config/ssl/ssl.crt/test.crt",
  "subject": "CN=test.example.com,O=Test Inc.,L=New York,ST=New York,C=US",
  "updatedBy": "admin",
  "version": 1
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	cert, err := s.Client.GetCertificate("/Common/test.crt")

	require.Nil(s.T(), err, "Error getting certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriFile, uriSslCert, "~Common~test.crt"), s.LastRequest.URL.Path, "Wrong uri to fetch certificate")
	s.requireReserializesTo(resp, cert, "Certificate should reserialize to itself")
}

func (s *SysTestSuite) TestDeleteCertificate() {
	err := s.Client.DeleteCertificate("/Common/test.crt")

	require.Nil(s.T(), err, "Error deleting certificate")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriFile, uriSslCert, "~Common~test.crt"), s.LastRequest.URL.Path, "Wrong uri to delete certificate")
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *SysTestSuite) TestKeys() {
	resp := `{
  "items": [
    {
      "name": "foo.key",
      "partition": "Common",
      "fullPath": "/Common/foo.key",
      "generation": 1,
      "checksum": "SHA1:1704:bbfccd4b6d1215b2f7feb92ba8ff336bdb22885f",
      "createTime": "2019-03-15T15:10:11Z",
      "createdBy": "admin",
      "curveName": "none",
      "keySize": 2048,
      "keyType": "rsa-private",
      "lastUpdateTime": "2019-03-15T15:10:11Z",
      "mode": 33184,
      "revision": 1,
      "securityType": "normal",
      "size": 1704,
      "sourcePath": "/config/ssl/ssl.key/foo.key",
      "updatedBy": "admin"
    },
    {
      "name": "bar.key",
      "partition": "Common",
      "fullPath": "/Common/bar.key",
      "generation": 1,
      "checksum": "SHA1:1704:b41c6b484421ad6229bc7e8ea623f0eacf5bc78b",
      "createTime": "2019-03-15T15:10:31Z",
      "createdBy": "admin",
      "curveName": "none",
      "keySize": 2048,
      "keyType": "rsa-private",
      "lastUpdateTime": "2019-03-15T15:10:31Z",
      "mode": 33184,
      "revision": 1,
      "securityType": "normal",
      "size": 1704,
      "sourcePath": "/config/ssl/ssl.key/bar.key",
      "updatedBy": "admin"
    }
  ]
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	keys, err := s.Client.Keys()

	require.Nil(s.T(), err, "Error loading keys")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFile, uriSslKey), s.LastRequest.URL.Path, "Wrong uri to fetch keys")
	assert.Equal(s.T(), 2, len(keys.Keys), "Wrong number of keys")
	assert.Equal(s.T(), "foo.key", keys.Keys[0].Name)
	assert.Equal(s.T(), "bar.key", keys.Keys[1].Name)
	s.requireReserializesTo(resp, keys, "Keys should reserialize to itself")
}

func (s *SysTestSuite) TestAddKey() {
	key := Key{
		Name:       "test",
		SourcePath: "file:///var/config/rest/downloads/test.key",
	}
	err := s.Client.AddKey(&key)

	require.Nil(s.T(), err, "Error adding key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriSys, uriFile, uriSslKey), s.LastRequest.URL.Path, "Wrong uri to create key")
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"test", "sourcePath":"file:///var/config/rest/downloads/test.key"}`, s.LastRequestBody)
}

func (s *SysTestSuite) TestGetKey() {
	resp := `{
  "name": "test.key",
  "partition": "Common",
  "fullPath": "/Common/test.key",
  "generation": 1,
  "checksum": "SHA1:1704:73c6766b89be06a464d2269a0b96b2cc0e6cc2f2",
  "createTime": "2019-03-15T15:09:46Z",
  "createdBy": "admin",
  "curveName": "none",
  "keySize": 2048,
  "keyType": "rsa-private",
  "lastUpdateTime": "2019-03-15T15:09:46Z",
  "mode": 33184,
  "revision": 1,
  "securityType": "normal",
  "size": 1704,
  "sourcePath": "/config/ssl/ssl.key/test.key",
  "updatedBy": "admin"
}`

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	key, err := s.Client.GetKey("/Common/test.key")

	require.Nil(s.T(), err, "Error getting key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriFile, uriSslKey, "~Common~test.key"), s.LastRequest.URL.Path, "Wrong uri to fetch key")
	s.requireReserializesTo(resp, key, "Key should reserialize to itself")
}

func (s *SysTestSuite) TestDeleteKey() {
	err := s.Client.DeleteKey("/Common/test.key")

	require.Nil(s.T(), err, "Error deleting key")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriSys, uriFile, uriSslKey, "~Common~test.key"), s.LastRequest.URL.Path, "Wrong uri to delete key")
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}
