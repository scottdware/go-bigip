package bigip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GTMTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *GTMTestSuite) SetupSuite() {
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

func (s *GTMTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *GTMTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestGtmSuite(t *testing.T) {
	suite.Run(t, new(GTMTestSuite))
}

func (s *GTMTestSuite) TestGTMWideIPsARecord() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(wideIPsSample())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	w, err := s.Client.GTMWideIPs(ARecord)

	// make sure we get wideIp's back
	assert.NotNil(w)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriWideIp, uriARecord), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", w.GTMWideIPs[0].Partition)
	assert.Equal("/Common/baseapp.domain.com", w.GTMWideIPs[0].FullPath)
	// ensure we can find our partition-based WideIP
	assert.Equal("test", w.GTMWideIPs[1].Partition)
	assert.Equal("/test/myapp.domain.com", w.GTMWideIPs[1].FullPath)

}

func (s *GTMTestSuite) TestGetGTMWideIP() {
	// ** Test Common (partition)

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		// get sample WideIP in Common partition
		w.Write(wideIPSample(false))
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	w, err := s.Client.GetGTMWideIP("baseapp.domain.com", ARecord)

	// make sure we get wideIp's back
	assert.NotNil(w)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriWideIp, uriARecord, "baseapp.domain.com"), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", w.Partition)
	assert.Equal("/Common/baseapp.domain.com", w.FullPath)

	// ** Test Partition

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		// get sample WideIP in test partition
		w.Write(wideIPSample(true))
	}

	w2, err := s.Client.GetGTMWideIP("/test/myapp.domain.com", ARecord)

	// make sure we get wideIp's back
	assert.NotNil(w2)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriWideIp, uriARecord, "~test~myapp.domain.com"), s.LastRequest.URL.Path)
	// ensure we can find our partition-based WideIP
	assert.Equal("test", w2.Partition)
	assert.Equal("/test/myapp.domain.com", w2.FullPath)
}
