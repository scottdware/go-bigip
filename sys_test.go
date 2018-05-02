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
