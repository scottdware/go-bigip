package bigip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"strings"
)

type LTMTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *LTMTestSuite) SetupSuite() {
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

func (s *LTMTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *LTMTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestLtmSuite(t *testing.T) {
	suite.Run(t, new(LTMTestSuite))
}

func (s *LTMTestSuite) TestIRules() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"items" :
			[
				{"name":"rule1","apiAnonymous":"rule1"},
				{"name":"rule2","apiAnonymous":"this\nis\nrule2"}
			]}`))
	}

	rules, err := s.Client.IRules()

	assert.Nil(s.T(), err, "Error loading rules")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriIRule), s.LastRequest.URL.Path, "Wrong uri to fetch rules")
	assert.Equal(s.T(), 2, len(rules.IRules), "Wrong number of rules")
	assert.Equal(s.T(), "rule1", rules.IRules[0].Name)
	assert.Equal(s.T(), "rule2", rules.IRules[1].Name)
	assert.Equal(s.T(), `this
is
rule2`, rules.IRules[1].Rule, "Multiline rule not unmarshalled")
}

func (s *LTMTestSuite) TestCreateIRule() {
	s.Client.CreateIRule("rule1", `when CLIENT_ACCEPTED { log local0. "test"}`)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriIRule), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"when CLIENT_ACCEPTED { log local0. \"test\"}"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyIRule() {
	s.Client.ModifyIRule("rule1", &IRule{Rule: "modified"})

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"modified"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteIRule() {
	s.Client.DeleteIRule("rule1")

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestModifyVirtualAddress() {
	d := &VirtualAddress{}
	s.Client.ModifyVirtualAddress("address1", d)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriVirtualAddress, "address1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestGetPolicies() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:ltm:policy:policycollectionstate",
  "selfLink": "https://localhost/mgmt/tm/ltm/policy?ver=11.5.1",
  "items": [
    {
      "kind": "tm:ltm:policy:policystate",
      "name": "policy1",
      "partition": "Common",
      "fullPath": "/Common/policy1",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/ltm/policy/~Common~policy1?ver=11.5.1",
      "controls": [
        "forwarding"
      ],
      "requires": [
        "http",
        "client-ssl"
      ],
      "strategy": "/Common/first-match",
      "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~policy1/rules?ver=11.5.1",
        "isSubcollection": true
      }
    },
    {
      "kind": "tm:ltm:policy:policystate",
      "name": "policy2",
      "partition": "Common",
      "fullPath": "/Common/policy2",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/ltm/policy/~Common~policy2?ver=11.5.1",
      "controls": [
        "classification"
      ],
      "requires": [
        "ssl-persistence"
      ],
      "strategy": "/Common/first-match",
      "rulesReference": {
        "link": "https://localhost/mgmt/tm/ltm/policy/~Common~policy2/rules?ver=11.5.1",
        "isSubcollection": true
      }
    }
    ]}`))
	}

	p, e := s.Client.Policies()

	assert.Nil(s.T(), e, "Fetching policy list should not return an error")
	assert.Equal(s.T(), 2, len(p.Policies), "Wrong number of policies returned")
	assert.Equal(s.T(), "policy1", p.Policies[0].Name)
	assert.Equal(s.T(), "Common", p.Policies[0].Partition)
	assert.Equal(s.T(), "/Common/first-match", p.Policies[0].Strategy)
	assert.EqualValues(s.T(), []string{"http", "client-ssl"}, p.Policies[0].Requires)
	assert.EqualValues(s.T(), []string{"forwarding"}, p.Policies[0].Controls)
}

func (s *LTMTestSuite) TestGetPolicy() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "rules") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:rulescollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:rulesstate",
			      "name": "rule1",
			      "fullPath": "rule1",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1?ver=11.5.1",
			      "ordinal": 0,
			      "actionsReference": {
				"link": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions?ver=11.5.1",
				"isSubcollection": true
			      },
			      "conditionsReference": {
				"link": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions?ver=11.5.1",
				"isSubcollection": true
			      }
			    }
			  ]
			}`))
		} else if strings.HasSuffix(r.URL.Path, "actions") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:actions:actionscollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:actions:actionsstate",
			      "name": "0",
			      "fullPath": "0",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/actions/0?ver=11.5.1",
			      "code": 0,
			      "forward": true,
			      "pool": "/Common/sorry_server",
			      "port": 0,
			      "request": true,
			      "select": true,
			      "status": 0,
			      "vlanId": 0
			    }
			  ]
			}`))
		} else if strings.HasSuffix(r.URL.Path, "conditions") {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:rules:conditions:conditionscollectionstate",
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions?ver=11.5.1",
			  "items": [
			    {
			      "kind": "tm:ltm:policy:rules:conditions:conditionsstate",
			      "name": "0",
			      "fullPath": "0",
			      "generation": 144344,
			      "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy/rules/rule1/conditions/0?ver=11.5.1",
			      "all": true,
			      "caseInsensitive": true,
			      "external": true,
			      "httpUri": true,
			      "index": 0,
			      "present": true,
			      "remote": true,
			      "request": true,
			      "startsWith": true,
			      "values": [
				"/foo"
			      ]
			    }]
			    }`))
		} else {
			w.Write([]byte(`{
			  "kind": "tm:ltm:policy:policystate",
			  "name": "my_policy",
			  "fullPath": "my_policy",
			  "generation": 144344,
			  "selfLink": "https://localhost/mgmt/tm/ltm/policy/my_policy?ver=11.5.1",
			  "controls": [
			    "forwarding"
			  ],
			  "requires": [
			    "http"
			  ],
			  "strategy": "/Common/first-match",
			  "rulesReference": {
			    "link": "https://localhost/mgmt/tm/ltm/policy/~Common~my_policy/rules?ver=11.5.1",
			    "isSubcollection": true
			  }
			}`))
		}
	}

	p, err := s.Client.GetPolicy("my_policy")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "my_policy", p.Name)
	assert.Equal(s.T(), 1, len(p.Rules), "Not enough rules")
	assert.Equal(s.T(), 1, len(p.Rules[0].Actions), "Not enough actions")
	assert.Equal(s.T(), "/Common/sorry_server", p.Rules[0].Actions[0].Pool)
	assert.Equal(s.T(), 1, len(p.Rules[0].Conditions), "Not enough conditions")
	assert.Equal(s.T(), []string{"/foo"}, p.Rules[0].Conditions[0].Values)
}

func (s *LTMTestSuite) TestGetNonExistentPolicy() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	p, err := s.Client.GetPolicy("asdf")

	assert.NotNil(s.T(), err, "nil error returned")
	assert.Nil(s.T(), p)
	assert.True(s.T(), strings.HasPrefix(err.Error(), "HTTP 404"), err.Error())
}

func (s *LTMTestSuite) TestCreatePolicy() {
	p := Policy{
		Name:     "test",
		Controls: []string{"forwarding"},
		Requires: []string{"http"},
		Strategy: "/Common/first-match",
		Rules: []PolicyRule{
			PolicyRule{
				Name: "rule1",
				Actions: []PolicyRuleAction{
					PolicyRuleAction{
						Forward: true,
						Pool:    "somepool",
					},
					PolicyRuleAction{},
				},
				Conditions: []PolicyRuleCondition{
					PolicyRuleCondition{
						CaseInsensitive: true,
						Values:          []string{"/foo", "/bar"},
					},
					PolicyRuleCondition{},
				},
			},
			PolicyRule{
				Name: "rule2",
			},
		},
	}

	s.Client.CreatePolicy(&p)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriPolicy), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `{"name":"test",
		"controls":["forwarding"],
		"requires":["http"],
		"strategy":"/Common/first-match",
		"rulesReference" : {"items":[
			{ "name":"rule1",
			  "ordinal":0,
			  "actionsReference":{ "items":[
			  	{
			  		"name":"0",
			  		"pool":"somepool",
			  		"forward":true
			  	},
			  	{
			  		"name":"1"
			  	}
			  ]},
			  "conditionsReference" : {"items":[
			  	{
			  		"name":"0",
			  		"caseInsensitive":true,
			  		"values":["/foo","/bar"]
			  	},
			  	{
			  		"name":"1"
			  	}
			  ]}
			},
			{ "name":"rule2", "ordinal":1, "actionsReference":{}, "conditionsReference":{}}]}}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestUpdatePolicy() {
	//TODO: test more stuff
	s.Client.UpdatePolicy("foo", &Policy{})

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo", uriLtm, uriPolicy), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestDeletePolicy() {
	s.Client.DeletePolicy("foo")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo", uriLtm, uriPolicy), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestCreateVitualAddress() {

	s.Client.CreateVirtualAddress("test-va", &VirtualAddress{Address: "10.10.10.10", ARP: true, AutoDelete: false})

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriLtm, uriVirtualAddress), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `
	{"name":"test-va",
	"arp":"enabled",
	"autoDelete":"false",
	"address" : "10.10.10.10",
	"enabled":"no",
	"floating":"disabled",
	"icmpEcho":"disabled",
	"inheritedTrafficGroup":"no",
	"routeAdvertisement":"disabled"}`, s.LastRequestBody)

}

func (s *LTMTestSuite) TestDeleteVitualAddress() {

	s.Client.DeleteVirtualAddress("test-va")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/test-va", uriLtm, uriVirtualAddress), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestModifyVirtualServer() {
	vs := &VirtualServer{
		Name: "test",
		Profiles: []Profile{
			Profile{Name: "/Common/tcp", Context: CONTEXT_CLIENT},
			Profile{Name: "/Common/tcp", Context: CONTEXT_SERVER}},
		//TODO: test more
	}

	s.Client.ModifyVirtualServer("test", vs)

	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/test", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.JSONEq(s.T(), `
	{"name":"test",
	"sourceAddressTranslation":{},
	"profiles":[
		{"name":"/Common/tcp","context":"clientside"},
		{"name":"/Common/tcp","context":"serverside"}
	]
	}`, s.LastRequestBody)

}

func (s *LTMTestSuite) TestCreateMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
		Interval:      15,
		Timeout:       5,
		SendString:    "GET /\r\n",
		ReceiveString: "200 OK",
	}

	s.Client.CreateMonitor(config.Name, config.ParentMonitor, config.Interval, config.Timeout, config.SendString, config.ReceiveString)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-web-monitor","defaultsFrom":"http","interval":15,"manualResume":"disabled","recv":"200 OK","reverse":"disabled","send":"GET /\\r\\n","timeout":5,"transparent":"disabled"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestAddMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
		Interval:      15,
		Timeout:       5,
		SendString:    "GET /\r\n",
		ReceiveString: "200 OK",
		Username:      "monitoring",
		Password:      "monitoring",
	}

	s.Client.AddMonitor(config)

	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor), s.LastRequest.URL.Path)
	assert.Equal(s.T(), `{"name":"test-web-monitor","defaultsFrom":"http","interval":15,"manualResume":"disabled","password":"monitoring","recv":"200 OK","reverse":"disabled","send":"GET /\\r\\n","timeout":5,"transparent":"disabled","username":"monitoring"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteMonitor() {
	config := &Monitor{
		Name:          "test-web-monitor",
		ParentMonitor: "http",
	}

	s.Client.DeleteMonitor(config.Name, config.ParentMonitor)

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s/%s", uriLtm, uriMonitor, config.ParentMonitor, config.Name), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestVirtualServerPolicies() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		  "kind": "tm:ltm:virtual:policies:policiescollectionstate",
		  "selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies?ver=11.5.1",
		  "policiesReference": {
		    "link": "https://localhost/mgmt/tm/ltm/virtual/foo/policies?ver=11.5.1",
		    "isSubcollection": true,
		    "items": [
		      {
			"kind": "tm:ltm:virtual:policies:policiesstate",
			"name": "policy1",
			"partition": "Common",
			"fullPath": "/Common/policy1",
			"generation": 1,
			"selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies/~Common~policy1?ver=11.5.1"
		      },
		      {
			"kind": "tm:ltm:virtual:policies:policiesstate",
			"name": "policy2",
			"partition": "Common",
			"fullPath": "/Common/policy2",
			"generation": 1,
			"selfLink": "https://localhost/mgmt/tm/ltm/virtual/foo/policies/~Common~policy2?ver=11.5.1"
		      }
		    ]
		  }
		}`))
	}

	p, err := s.Client.VirtualServerPolicyNames("foo")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/foo/policies", uriLtm, uriVirtual), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "/Common/policy1", p[0])
	assert.Equal(s.T(), "/Common/policy2", p[1])
}

func (s *LTMTestSuite) TestInternalDataGroups() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
		  "kind": "tm:ltm:data-group:internal:internalcollectionstate",
		  "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal?ver=12.1.2",
		  "items": [
		    {
		      "kind": "tm:ltm:data-group:internal:internalstate",
		      "name": "some_data_group",
		      "partition": "Common",
		      "fullPath": "/Common/api.tv2.dk_host_pool_map",
		      "generation": 2552,
		      "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal/~Common~api.tv2.dk_host_pool_map?ver=12.1.2",
		      "type": "string",
		      "records": [
		        {
		          "name": "jens.medister.api.tv2.dk-hest",
		          "data": "pool-medister"
		        }
		      ]
		    },
		    {
		      "kind": "tm:ltm:data-group:internal:internalstate",
		      "name": "jenkins_whitelisted_paths",
		      "partition": "Common",
		      "fullPath": "/Common/jenkins_whitelisted_paths",
		      "generation": 41,
		      "selfLink": "https://localhost/mgmt/tm/ltm/data-group/internal/~Common~jenkins_whitelisted_paths?ver=12.1.2",
		      "type": "string",
		      "records": [
		        {
		          "name": "/medister",
		          "data": "1"
		        }
		      ]
		    }
		  ]
		}`))
	}

	g, err := s.Client.InternalDataGroups()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "GET", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s/%s", uriLtm, uriDatagroup, uriInternal), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "some_data_group", g.DataGroups[0].Name)
	assert.Equal(s.T(), "jenkins_whitelisted_paths", g.DataGroups[1].Name)
}
