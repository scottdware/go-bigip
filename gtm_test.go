package bigip

import (
	"encoding/json"
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

// ********************************************************************************************************************
// *************************************************                  *************************************************
// *************************************************   GTM WideIP A   *************************************************
// *************************************************                  *************************************************
// ********************************************************************************************************************

func (s *GTMTestSuite) TestGTMWideIPsARecord() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(wideIPSamples())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	w, err := s.Client.GetGTMWideIPs(ARecord)

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

func (s *GTMTestSuite) TestGetGTMWideIPARecord() {
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

func (s *GTMTestSuite) TestAddGTMWideIPARecord() {
	config := &GTMWideIP{
		Name:      "baseapp.domain.com",
		Partition: "Common",
	}

	s.Client.AddGTMWideIP(config, ARecord)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriWideIp, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com","partition":"Common"}`, s.LastRequestBody)

}

func (s *GTMTestSuite) TestAddGTMWideIPAdvancedARecord() {
	config := &GTMWideIP{}
	if err := json.Unmarshal(wideIPSample(false), &config); err != nil {
		panic(err)
	}

	// make sure our post works
	err := s.Client.AddGTMWideIP(config, ARecord)
	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test no error on post
	assert.Nil(err)
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriWideIp, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(wideIPReturn(false), s.LastRequestBody)

}

func (s *GTMTestSuite) TestDeleteGTMWideIPARecord() {
	fullPath := "/Common/baseapp.domain.com"

	s.Client.DeleteGTMWideIP(fullPath, ARecord)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriWideIp, uriARecord, "~Common~baseapp.domain.com"), s.LastRequest.URL.Path)
}

func (s *GTMTestSuite) TestModifyGTMWideIPARecord() {
	config := &GTMWideIP{
		Name:      "baseapp.domain.com",
		Partition: "Common",
	}

	s.Client.AddGTMWideIP(config, ARecord)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriWideIp, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com","partition":"Common"}`, s.LastRequestBody)

	configUpdate := &GTMWideIP{
		Name:      "baseapp.domain.com",
		Partition: "Common",
		Pools: &[]GTMWideIPPool{
			{
				Partition: "Common",
				Name:      "baseapp.domain.com_pool",
			},
		},
	}

	fullPath := "/Common/baseapp.domain.com"

	s.Client.ModifyGTMWideIP(fullPath, configUpdate, ARecord)

	// Test we put
	assert.Equal("PUT", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriWideIp, uriARecord, "~Common~baseapp.domain.com"), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com","partition":"Common","pools":[{"name":"baseapp.domain.com_pool","partition":"Common","nameReference":{}}]}`, s.LastRequestBody)

}

// ********************************************************************************************************************
// *************************************************                ***************************************************
// *************************************************   GTM Pool A   ***************************************************
// *************************************************                ***************************************************
// ********************************************************************************************************************

func (s *GTMTestSuite) TestGTMPoolARecord() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(poolASamples())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	w, err := s.Client.GetGTMAPools()

	// make sure we get wideIp's back
	assert.NotNil(w)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriPool, uriARecord), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", w.GTMAPools[0].Partition)
	assert.Equal("/Common/baseapp.domain.com_pool", w.GTMAPools[0].FullPath)
	// ensure we can find our partition-based WideIP
	assert.Equal("test", w.GTMAPools[1].Partition)
	assert.Equal("/test/myapp.domain.com_pool", w.GTMAPools[1].FullPath)

}

func (s *GTMTestSuite) TestGetGTMPoolARecord() {
	// ** Test Common (partition)

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		// get sample Pool in Common partition
		w.Write(poolASample(false))
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	p, err := s.Client.GetGTMAPool("baseapp.domain.com")

	// make sure we get wideIp's back
	assert.NotNil(p)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, "baseapp.domain.com"), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", p.Partition)
	assert.Equal("/Common/baseapp.domain.com_pool", p.FullPath)

	// ** Test Partition

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		// get sample Pool in test partition
		w.Write(poolASample(true))
	}

	p2, err := s.Client.GetGTMAPool("/test/myapp.domain.com")

	// make sure we get wideIp's back
	assert.NotNil(p2)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, "~test~myapp.domain.com"), s.LastRequest.URL.Path)
	// ensure we can find our partition-based Pool
	assert.Equal("test", p2.Partition)
	assert.Equal("/test/myapp.domain.com_pool", p2.FullPath)
}

func (s *GTMTestSuite) TestAddGTMPoolARecord() {
	config := &GTMAPool{
		Name:      "baseapp.domain.com_pool",
		Partition: "Common",
	}

	s.Client.AddGTMAPool(config)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriPool, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com_pool","partition":"Common","MembersReference":{}}`, s.LastRequestBody)

}

func (s *GTMTestSuite) TestAddGTMPoolAdvancedARecord() {
	config := &GTMAPool{}
	if err := json.Unmarshal(poolASample(false), &config); err != nil {
		panic(err)
	}

	// make sure our post works
	err := s.Client.AddGTMAPool(config)
	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test no error on post
	assert.Nil(err)
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriPool, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(poolAReturn(false), s.LastRequestBody)

	// Try partition
	config2 := &GTMAPool{}
	if err := json.Unmarshal(poolASample(true), &config2); err != nil {
		panic(err)
	}
	// make sure our post works
	err2 := s.Client.AddGTMAPool(config2)
	// Test no error on post
	assert.Nil(err2)
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriPool, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(poolAReturn(true), s.LastRequestBody)
}

func (s *GTMTestSuite) TestDeleteGTMPoolARecord() {
	fullPath := "/Common/baseapp.domain.com"

	s.Client.DeleteGTMPool(fullPath, ARecord)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, "~Common~baseapp.domain.com"), s.LastRequest.URL.Path)
}

func (s *GTMTestSuite) TestModifyGTMPoolARecord() {
	config := &GTMAPool{
		Name:      "baseapp.domain.com",
		Partition: "Common",
	}

	s.Client.AddGTMAPool(config)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriGtm, uriPool, uriARecord), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com","partition":"Common","MembersReference":{}}`, s.LastRequestBody)

	configUpdate := &GTMAPool{
		Name:      "baseapp.domain.com_pool",
		Partition: "Common",
		Disabled:  true,
		Enabled:   false,
	}

	fullPath := "/Common/baseapp.domain.com_pool"

	s.Client.ModifyGTMAPool(fullPath, configUpdate)

	// Test we put
	assert.Equal("PUT", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, "~Common~baseapp.domain.com_pool"), s.LastRequest.URL.Path)
	// See if we get back the object we expect
	assert.Equal(`{"name":"baseapp.domain.com_pool","partition":"Common","disabled":true,"MembersReference":{}}`, s.LastRequestBody)

}

// ********************************************************************************************************************
// *****************************************                        ***************************************************
// *****************************************   GTM A Pool Members   ***************************************************
// *****************************************                        ***************************************************
// ********************************************************************************************************************

func (s *GTMTestSuite) TestGetGTMAPoolMembers() {
	fullPathAPool := "/Common/baseapp.domain.com_pool"
	fullPathAPoolAPI := "~Common~baseapp.domain.com_pool"

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(poolAMemberSamples())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	m, err := s.Client.GetGTMAPoolMembers(fullPathAPool)

	// make sure we get wideIp's back
	assert.NotNil(m)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, fullPathAPoolAPI, uriPoolMembers), s.LastRequest.URL.Path)
	// ensure we can find our common WideIp
	assert.Equal("Common", m.GTMAPoolMembers[0].Partition)
	assert.Equal("someltm:/Common", m.GTMAPoolMembers[0].SubPath)
	assert.Equal("/Common/someltm:/Common/baseapp_80_vs", m.GTMAPoolMembers[0].FullPath)

}

func (s *GTMTestSuite) TestGetGTMAPoolMember() {
	fullPathAPool := "/Common/baseapp.domain.com_pool"
	fullPathAPoolAPI := "~Common~baseapp.domain.com_pool"
	memberPath := "/Common/baseapp_80_vs"
	memberPathAPI := "~Common~baseapp_80_vs"
	serverPath := "/Common/someltm"
	serverPathAPI := "~Common~someltm"
	theCrazyPoolMemberFullPath := fmt.Sprintf("%s:%s", serverPath, memberPath)
	theCrazyPoolMemberFullPathAPI := fmt.Sprintf("%s:%s", serverPathAPI, memberPathAPI)

	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write(poolAMemberSample())
	}

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	m, err := s.Client.GetGTMAPoolMember(fullPathAPool, serverPath, memberPath)

	// make sure we get wideIp's back
	assert.NotNil(m)
	assert.Nil(err)

	// see that we talked to the gtm/wideip/a endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, fullPathAPoolAPI, uriPoolMembers, theCrazyPoolMemberFullPathAPI), s.LastRequest.URL.Path)
	assert.Equal("baseapp_80_vs", m.Name)
	assert.Equal("Common", m.Partition)
	assert.Equal("someltm:/Common", m.SubPath)
	assert.Equal(theCrazyPoolMemberFullPath, m.FullPath)
}

func (s *GTMTestSuite) TestCreateGTMAPoolMember() {
	fullPathAPool := "/Common/baseapp.domain.com_pool"
	fullPathAPoolAPI := "~Common~baseapp.domain.com_pool"
	memberPath := "/Common/baseapp_80_vs"
	serverPath := "/Common/someltm"

	s.Client.CreateGTMAPoolMember(fullPathAPool, serverPath, memberPath)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("POST", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, fullPathAPoolAPI, uriPoolMembers), s.LastRequest.URL.Path)
	assert.Equal(`{"name":"/Common/someltm:/Common/baseapp_80_vs"}`, s.LastRequestBody)

}

func (s *GTMTestSuite) TestDeleteGTMAPoolMember() {
	fullPathAPool := "/Common/baseapp.domain.com_pool"
	fullPathAPoolAPI := "~Common~baseapp.domain.com_pool"
	memberPath := "/Common/baseapp_80_vs"
	memberPathAPI := "~Common~baseapp_80_vs"
	serverPath := "/Common/someltm"
	serverPathAPI := "~Common~someltm"
	theCrazyPoolMemberFullPathAPI := fmt.Sprintf("%s:%s", serverPathAPI, memberPathAPI)

	s.Client.DeleteGTMAPoolMember(fullPathAPool, serverPath, memberPath)

	// so we don't have to pass  s.T() as first argument every time in Assert
	assert := assert.New(s.T())
	// Test we posted
	assert.Equal("DELETE", s.LastRequest.Method)
	// See that we actually posted to our endpoint
	assert.Equal(fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s/%s/%s", uriGtm, uriPool, uriARecord, fullPathAPoolAPI, uriPoolMembers, theCrazyPoolMemberFullPathAPI), s.LastRequest.URL.Path)

}

// ********************************************************************************************************************
// **********************************************               *******************************************************
// **********************************************   Test Data   *******************************************************
// **********************************************               *******************************************************
// ********************************************************************************************************************

func wideIPSamples() []byte {
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
										"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1"
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

func wideIPSample(usePartition bool) []byte {
	if usePartition {
		return []byte(`{
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
		}`)
	}

	return []byte(`{
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
	}`)
}

func wideIPReturn(usePartiion bool) string {
	if usePartiion {
		return `{"name":"myapp.domain.com","partition":"test","fullPath":"/test/baseapp.domain.com","generation":2,"enabled":true,"failureRcode":"noerror","failureRcodeResponse":"disabled","minimalResponse":"enabled","persistCidrIpv4":32,"persistCidrIpv6":128,"persistence":"disabled","poolLbMode":"topology","ttlPersistence":3600,"pools":[{"name":"myapp.domain.com_pool","partition":"test","ratio":1,"nameReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool_int_pool?ver=12.1.1"}}]}`
	}

	return `{"name":"baseapp.domain.com","partition":"Common","fullPath":"/Common/baseapp.domain.com","generation":2,"enabled":true,"failureRcode":"noerror","failureRcodeResponse":"disabled","minimalResponse":"enabled","persistCidrIpv4":32,"persistCidrIpv6":128,"persistence":"disabled","poolLbMode":"topology","ttlPersistence":3600,"pools":[{"name":"baseapp.domain.com_pool","partition":"Common","ratio":1,"nameReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool_int_pool?ver=12.1.1"}}]}`
}

func poolASamples() []byte {
	return []byte(
		`{
			"kind": "tm:gtm:pool:a:acollectionstate",
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a?ver=12.1.1",
			"items": [
					{
							"kind": "tm:gtm:pool:a:astate",
							"name": "baseapp.domain.com_pool",
							"partition": "Common",
							"fullPath": "/Common/baseapp.domain.com_pool",
							"generation": 2,
							"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1",
							"alternateMode": "round-robin",
							"dynamicRatio": "disabled",
							"enabled": true,
							"fallbackIp": "any",
							"fallbackMode": "return-to-dns",
							"limitMaxBps": 0,
							"limitMaxBpsStatus": "disabled",
							"limitMaxConnections": 0,
							"limitMaxConnectionsStatus": "disabled",
							"limitMaxPps": 0,
							"limitMaxPpsStatus": "disabled",
							"loadBalancingMode": "round-robin",
							"manualResume": "disabled",
							"maxAnswersReturned": 1,
							"monitor": "default",
							"qosHitRatio": 5,
							"qosHops": 0,
							"qosKilobytesSecond": 3,
							"qosLcs": 30,
							"qosPacketRate": 1,
							"qosRtt": 50,
							"qosTopology": 0,
							"qosVsCapacity": 0,
							"qosVsScore": 0,
							"ttl": 30,
							"verifyMemberAvailability": "enabled",
							"membersReference": {
									"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
									"isSubcollection": true
							}
					},
					{
            "kind": "tm:gtm:pool:a:astate",
            "name": "myapp.domain.com_pool",
            "partition": "test",
            "fullPath": "/test/myapp.domain.com_pool",
            "generation": 182,
            "selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1",
            "alternateMode": "round-robin",
            "dynamicRatio": "disabled",
            "enabled": true,
            "fallbackIp": "any",
            "fallbackMode": "return-to-dns",
            "limitMaxBps": 0,
            "limitMaxBpsStatus": "disabled",
            "limitMaxConnections": 0,
            "limitMaxConnectionsStatus": "disabled",
            "limitMaxPps": 0,
            "limitMaxPpsStatus": "disabled",
            "loadBalancingMode": "round-robin",
            "manualResume": "disabled",
            "maxAnswersReturned": 1,
            "monitor": "default",
            "qosHitRatio": 5,
            "qosHops": 0,
            "qosKilobytesSecond": 3,
            "qosLcs": 30,
            "qosPacketRate": 1,
            "qosRtt": 50,
            "qosTopology": 0,
            "qosVsCapacity": 0,
            "qosVsScore": 0,
            "ttl": 30,
            "verifyMemberAvailability": "enabled",
            "membersReference": {
                "link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1",
                "isSubcollection": true
            }
          }
      ]
    }`)
}

func poolASample(usePartition bool) []byte {
	if usePartition {
		return []byte(`{
			"kind": "tm:gtm:pool:a:astate",
			"name": "myapp.domain.com_pool",
			"partition": "test",
			"fullPath": "/test/myapp.domain.com_pool",
			"generation": 182,
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool?ver=12.1.1",
			"alternateMode": "round-robin",
			"dynamicRatio": "disabled",
			"enabled": true,
			"fallbackIp": "any",
			"fallbackMode": "return-to-dns",
			"limitMaxBps": 0,
			"limitMaxBpsStatus": "disabled",
			"limitMaxConnections": 0,
			"limitMaxConnectionsStatus": "disabled",
			"limitMaxPps": 0,
			"limitMaxPpsStatus": "disabled",
			"loadBalancingMode": "round-robin",
			"manualResume": "disabled",
			"maxAnswersReturned": 1,
			"monitor": "default",
			"qosHitRatio": 5,
			"qosHops": 0,
			"qosKilobytesSecond": 3,
			"qosLcs": 30,
			"qosPacketRate": 1,
			"qosRtt": 50,
			"qosTopology": 0,
			"qosVsCapacity": 0,
			"qosVsScore": 0,
			"ttl": 30,
			"verifyMemberAvailability": "enabled",
			"membersReference": {
					"link": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1",
					"isSubcollection": true
			}
		}`)
	}

	return []byte(`{
		"kind": "tm:gtm:pool:a:astate",
		"name": "baseapp.domain.com_pool",
		"partition": "Common",
		"fullPath": "/Common/baseapp.domain.com_pool",
		"generation": 2,
		"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool?ver=12.1.1",
		"alternateMode": "round-robin",
		"dynamicRatio": "disabled",
		"enabled": true,
		"fallbackIp": "any",
		"fallbackMode": "return-to-dns",
		"limitMaxBps": 0,
		"limitMaxBpsStatus": "disabled",
		"limitMaxConnections": 0,
		"limitMaxConnectionsStatus": "disabled",
		"limitMaxPps": 0,
		"limitMaxPpsStatus": "disabled",
		"loadBalancingMode": "round-robin",
		"manualResume": "disabled",
		"maxAnswersReturned": 1,
		"monitor": "default",
		"qosHitRatio": 5,
		"qosHops": 0,
		"qosKilobytesSecond": 3,
		"qosLcs": 30,
		"qosPacketRate": 1,
		"qosRtt": 50,
		"qosTopology": 0,
		"qosVsCapacity": 0,
		"qosVsScore": 0,
		"ttl": 30,
		"verifyMemberAvailability": "enabled",
		"membersReference": {
				"link": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
				"isSubcollection": true
		}
	}`)
}

func poolAReturn(usePartiion bool) string {
	if usePartiion {
		return `{"name":"myapp.domain.com_pool","partition":"test","fullPath":"/test/myapp.domain.com_pool","generation":182,"dynamicRatio":"disabled","enabled":true,"fallbackIp":"any","fallbackMode":"return-to-dns","limitMaxBpsStatus":"disabled","limitMaxConnectionsStatus":"disabled","limitMaxPpsStatus":"disabled","loadBalancingMode":"round-robin","manualResume":"disabled","maxAnswersReturned":1,"monitor":"default","qosHitRatio":5,"qosKilobytesSecond":3,"qosLcs":30,"qosPacketRate":1,"qosRtt":50,"ttl":30,"verifyMemberAvailability":"enabled","MembersReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members?ver=12.1.1","isSubcollection":true}}`
	}

	return `{"name":"baseapp.domain.com_pool","partition":"Common","fullPath":"/Common/baseapp.domain.com_pool","generation":2,"dynamicRatio":"disabled","enabled":true,"fallbackIp":"any","fallbackMode":"return-to-dns","limitMaxBpsStatus":"disabled","limitMaxConnectionsStatus":"disabled","limitMaxPpsStatus":"disabled","loadBalancingMode":"round-robin","manualResume":"disabled","maxAnswersReturned":1,"monitor":"default","qosHitRatio":5,"qosKilobytesSecond":3,"qosLcs":30,"qosPacketRate":1,"qosRtt":50,"ttl":30,"verifyMemberAvailability":"enabled","MembersReference":{"link":"https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1","isSubcollection":true}}`
}

func poolAMemberSamples() []byte {
	return []byte(
		`{
			"kind": "tm:gtm:pool:a:members:memberscollectionstate",
			"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members?ver=12.1.1",
			"items": [
				{
						"kind": "tm:gtm:pool:a:members:membersstate",
						"name": "baseapp_80_vs",
						"partition": "Common",
						"subPath": "someltm:/Common",
						"fullPath": "/Common/someltm:/Common/baseapp_80_vs",
						"generation": 197,
						"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members/~Common~someltm:~Common~baseapp_80_vs?ver=12.1.1",
						"enabled": true,
						"limitMaxBps": 0,
						"limitMaxBpsStatus": "disabled",
						"limitMaxConnections": 0,
						"limitMaxConnectionsStatus": "disabled",
						"limitMaxPps": 0,
						"limitMaxPpsStatus": "disabled",
						"memberOrder": 0,
						"monitor": "default",
						"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "baseapp_443_vs",
					"partition": "Common",
					"subPath": "someltm:/Common",
					"fullPath": "/Common/someltm:/Common/baseapp_443_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members/~Common~someltm:~Common~baseapp_443_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "myapp_80_vs",
					"partition": "Common",
					"subPath": "someltm:/test",
					"fullPath": "/Common/someltm:/test/myapp_80_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members/~Common~someltm:~test~myapp_80_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				},
				{
					"kind": "tm:gtm:pool:a:members:membersstate",
					"name": "myapp_443_vs",
					"partition": "Common",
					"subPath": "someltm:/test",
					"fullPath": "/Common/someltm:/test/myapp_443_vs",
					"generation": 197,
					"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~test~myapp.domain.com_pool/members/~Common~someltm:~test~myapp_443_vs?ver=12.1.1",
					"enabled": true,
					"limitMaxBps": 0,
					"limitMaxBpsStatus": "disabled",
					"limitMaxConnections": 0,
					"limitMaxConnectionsStatus": "disabled",
					"limitMaxPps": 0,
					"limitMaxPpsStatus": "disabled",
					"memberOrder": 0,
					"monitor": "default",
					"ratio": 1
				}
			]
		}`)
}

func poolAMemberSample() []byte {
	return []byte(
		`{
				"kind": "tm:gtm:pool:a:members:membersstate",
				"name": "baseapp_80_vs",
				"partition": "Common",
				"subPath": "someltm:/Common",
				"fullPath": "/Common/someltm:/Common/baseapp_80_vs",
				"generation": 197,
				"selfLink": "https://localhost/mgmt/tm/gtm/pool/a/~Common~baseapp.domain.com_pool/members/~Common~someltm:~Common~baseapp_80_vs?ver=12.1.1",
				"enabled": true,
				"limitMaxBps": 0,
				"limitMaxBpsStatus": "disabled",
				"limitMaxConnections": 0,
				"limitMaxConnectionsStatus": "disabled",
				"limitMaxPps": 0,
				"limitMaxPpsStatus": "disabled",
				"memberOrder": 0,
				"monitor": "default",
				"ratio": 1
		}`)
}
