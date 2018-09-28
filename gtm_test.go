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

func wideIPSample() []byte {
	return []byte(`{
		"kind": "tm:gtm:wideip:a:acollectionstate",
		"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a?ver=12.1.1",
		"items": [
			{
				"kind": "tm:gtm:wideip:a:astate",
				"name": "baseapp.domain.com",
				"partition": "Common",
				"fullPath": "/Common/baseapp.domain.com",
				"generation": 2,
				"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~Common~baseapp.domain.com?ver=12.1.1",
				"enabled": true,
				"failureRcode": "noerror",
				"failureRcodeResponse": "disabled",
				"failureRcodeTtl": 0,
				"lastResortPool": "",
				"minimalResponse": "enabled",
				"persistCidrIpv4": 32,
				"persistCidrIpv6": 128,
				"persistence": "disabled",
				"poolLbMode": "topology",
				"ttlPersistence": 3600,
				"pools": [
						{
								"name": "baseapp.domain.com_pool",
								"partition": "Common",
								"order": 0,
								"ratio": 1,
								"nameReference": {
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool_int_pool?ver=12.1.1"
								}
						}
				]
			},
			{
				"kind": "tm:gtm:wideip:a:astate",
				"name": "myapp.domain.com",
				"partition": "test",
				"fullPath": "/test/myapp.domain.com",
				"generation": 35,
				"selfLink": "https://localhost/mgmt/tm/gtm/wideip/a/~test~myapp.domain.com?ver=12.1.1",
				"enabled": true,
				"failureRcode": "noerror",
				"failureRcodeResponse": "disabled",
				"failureRcodeTtl": 0,
				"lastResortPool": "",
				"minimalResponse": "enabled",
				"persistCidrIpv4": 32,
				"persistCidrIpv6": 128,
				"persistence": "disabled",
				"poolLbMode": "round-robin",
				"ttlPersistence": 3600,
				"pools": [
						{
								"name": "myapp.domain.com.com_pool",
								"partition": "test",
								"order": 0,
								"ratio": 1,
								"nameReference": {
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1"
								}
						}
				]
			}
		]
	}`)
}

func (s *GTMTestSuite) TestGTMWideIPsARecord() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(wideIPSample())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	w, err := s.Client.GTMWideIPs(ARecord)
	// make sure we get wideIp's back
	assert.NotNil(w)
	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriWideIp, uriARecord), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", w.GTMWideIPs[0].Partition)
	assert.Equal("/Common/baseapp.domain.com", w.GTMWideIPs[0].FullPath)
	// ensure we can find our partition-based WideIP
	assert.Equal("test", w.GTMWideIPs[1].Partition)
	assert.Equal("/test/myapp.domain.com", w.GTMWideIPs[1].FullPath)
	assert.Nil(err)

}
