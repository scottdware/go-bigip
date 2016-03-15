package bigip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
)

type LTMTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter)
}

func (s *LTMTestSuite) SetupSuite() {
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		s.LastRequestBody = string(body)
		s.LastRequest = r
		if s.ResponseFunc != nil {
			s.ResponseFunc(w)
		}
	}))

	s.Client = NewSession(s.Server.URL, "", "")
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
	s.ResponseFunc = func(w http.ResponseWriter) {
		w.Write([]byte(`{"items" :
			[
				{"name":"rule1","apiAnonymous":"rule1"},
				{"name":"rule2","apiAnonymous":"this\nis\nrule2"}
			]}`))
	}

	rules, err := s.Client.IRules()

	assert.Nil(s.T(), err, "Error loading rules")
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s", uriIRule), s.LastRequest.URL.Path, "Wrong uri to fetch rules")
	assert.Equal(s.T(), 2, len(rules.IRules), "Wrong number of rules")
	assert.Equal(s.T(), "rule1", rules.IRules[0].Name)
	assert.Equal(s.T(), "rule2", rules.IRules[1].Name)
	assert.Equal(s.T(), `this
is
rule2`, rules.IRules[1].Rule, "Multiline rule not unmarshalled")
}

func (s *LTMTestSuite) TestCreateIRule() {
	s.Client.CreateIRule("rule1", `when CLIENT_ACCEPTED { log local0. "test"}`)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s", uriIRule), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"when CLIENT_ACCEPTED { log local0. \"test\"}"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyIRule() {
	s.Client.ModifyIRule("rule1", &IRule{Rule: "modified"})

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"modified"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDeleteIRule() {
	s.Client.DeleteIRule("rule1")

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriIRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestModifyVirtualAddress() {
	d := &VirtualAddress{}
	s.Client.ModifyVirtualAddress("address1", d)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriVirtualAddress, "address1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
}

func (s *LTMTestSuite) TestGetPolicies() {
	s.ResponseFunc = func(w http.ResponseWriter) {
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
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s", uriPolicy), s.LastRequest.URL.Path)
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
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/foo", uriPolicy), s.LastRequest.URL.Path)
}

func (s *LTMTestSuite) TestDeletePolicy() {
	s.Client.DeletePolicy("foo")

	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/foo", uriPolicy), s.LastRequest.URL.Path)
}
