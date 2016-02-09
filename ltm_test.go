package bigip

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"

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
		s.LastRequest = r;
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
	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s", uriRule), s.LastRequest.URL.Path, "Wrong uri to fetch rules")
	assert.Equal(s.T(), 2, len(rules.IRules), "Wrong number of rules")
	assert.Equal(s.T(), "rule1", rules.IRules[0].Name)
	assert.Equal(s.T(), "rule2", rules.IRules[1].Name)
	assert.Equal(s.T(), `this
is
rule2`, rules.IRules[1].Rule, "Multiline rule not unmarshalled")
}

func (s *LTMTestSuite) TestCreateIRule() {
	s.Client.CreateIRule("rule1", `when CLIENT_ACCEPTED { log local0. "test"}`)

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s", uriRule), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "POST", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"when CLIENT_ACCEPTED { log local0. \"test\"}"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestModifyIRule() {
	s.Client.ModifyIRule("rule1", &IRule{Rule:"modified"})

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "PUT", s.LastRequest.Method)
	assert.JSONEq(s.T(), `{"name":"rule1","apiAnonymous":"modified"}`, s.LastRequestBody)
}

func (s *LTMTestSuite) TestDelteIRule() {
	s.Client.DeleteIRule("rule1")

	assert.Equal(s.T(), fmt.Sprintf("/mgmt/tm/%s/%s", uriRule, "rule1"), s.LastRequest.URL.Path)
	assert.Equal(s.T(), "DELETE", s.LastRequest.Method)
}