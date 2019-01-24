package bigip

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NetTestSuite struct {
	suite.Suite
	Client          *BigIP
	Server          *httptest.Server
	LastRequest     *http.Request
	LastRequestBody string
	ResponseFunc    func(http.ResponseWriter, *http.Request)
}

func (s *NetTestSuite) SetupSuite() {
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

func (s *NetTestSuite) TearDownSuite() {
	s.Server.Close()
}

func (s *NetTestSuite) SetupTest() {
	s.ResponseFunc = nil
	s.LastRequest = nil
}

func TestNetSuite(t *testing.T) {
	suite.Run(t, new(NetTestSuite))
}

func (s *NetTestSuite) requireReserializesTo(expected string, actual interface{}, message string) {
	b, err := json.Marshal(actual)
	s.Require().Nil(err, message)

	s.Require().JSONEq(expected, string(b), message)
}

func (s *NetTestSuite) TestGetInterfaces() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:interface:interfacecollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/interface?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:interface:interfacestate",
      "name": "1.1",
      "fullPath": "1.1",
      "generation": 53,
      "selfLink": "https://localhost/mgmt/tm/net/interface/1.1?ver=11.5.1",
      "bundle": "not-supported",
      "enabled": true,
      "forceGigabitFiber": "disabled",
      "ifIndex": 48,
      "lldpAdmin": "txonly",
      "lldpTlvmap": 130943,
      "macAddress": "00:00:00:00:00:00",
      "mediaFixed": "10000T-FD",
      "mediaMax": "10000T-FD",
      "mediaSfp": "auto",
      "mtu": 9198,
      "preferPort": "sfp",
      "sflow": {
        "pollInterval": 0,
        "pollIntervalGlobal": "yes"
      },
      "stp": "enabled",
      "stpAutoEdgePort": "enabled",
      "stpEdgePort": "true",
      "stpLinkType": "auto"
    }]}`))
	}

	i, err := s.Client.Interfaces()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/interface", "")
	assert.Equal(s.T(), 1, len(i.Interfaces))
	assert.Equal(s.T(), "1.1", i.Interfaces[0].Name)
}

func (s *NetTestSuite) TestAddInterfaceToVLan() {
	err := s.Client.AddInterfaceToVlan("vlan-name", "iface-name", false)

	assert.Nil(s.T(), err)
	assertRestCall(s,
		"PUT",
		"/mgmt/tm/net/vlan/vlan-name/interfaces",
		`{"name":"iface-name", "untagged":true}`)
}

func (s *NetTestSuite) TestSelfIPs() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:self:selfcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/self?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:self:selfstate",
      "name": "0.0.0.0",
      "partition": "Common",
      "fullPath": "/Common/0.0.0.0",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/self/~Common~0.0.0.0?ver=11.5.1",
      "address": "0.0.0.0/20",
      "floating": "disabled",
      "inheritedTrafficGroup": "false",
      "trafficGroup": "/Common/traffic-group-local-only",
      "unit": 0,
      "vlan": "/Common/vlan-name",
      "allowService": "all"
    }]}`))
	}

	ips, err := s.Client.SelfIPs()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/self", "")
	assert.Equal(s.T(), 1, len(ips.SelfIPs))
	assert.Equal(s.T(), "0.0.0.0", ips.SelfIPs[0].Name)

}

func (s *NetTestSuite) TestCreateSelfIP() {
	err := s.Client.CreateSelfIP("0.0.0.0", "0.0.0.0/20", "vlan")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/self", `{"name":"0.0.0.0","address":"0.0.0.0/20", "vlan":"vlan"}`)
}

func (s *NetTestSuite) TestDeleteSelfIP() {
	err := s.Client.DeleteSelfIP("0.0.0.0")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/self/0.0.0.0", "")
}

func (s *NetTestSuite) TestModifySelfIP() {
	ip := &SelfIP{Address: "0.0.0.0/24"}

	err := s.Client.ModifySelfIP("0.0.0.0", ip)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/self/0.0.0.0", "")
}

func (s *NetTestSuite) TestTrunks() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:self:selfcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/trunk?ver=11.5.1",
  "items": [
  {
  	"name" : "trunk-name",
  	"id" : 1
  }
  ]}`))
	}

	trunks, err := s.Client.Trunks()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/trunk", "")
	assert.Equal(s.T(), 1, len(trunks.Trunks))
}

func (s *NetTestSuite) TestCreateTrunk() {
	err := s.Client.CreateTrunk("name", "ifaces", true)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/trunk", `{"name":"name", "interfaces":["ifaces"], "lacp":"enabled"}`)
}

func (s *NetTestSuite) TestDeleteTrunk() {
	err := s.Client.DeleteTrunk("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/trunk/name", "")
}

func (s *NetTestSuite) TestModifyTrunk() {
	trunk := &Trunk{Name: "name", LACP: "enabled"}

	err := s.Client.ModifyTrunk("name", trunk)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/trunk/name", `{"name":"name", "lacp":"enabled"}`)
}

func (s *NetTestSuite) TestVlans() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:vlan:vlancollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/vlan?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:vlan:vlanstate",
      "name": "vlan",
      "partition": "Common",
      "fullPath": "/Common/vlan",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/vlan/~Common~vlan?ver=11.5.1",
      "autoLasthop": "default",
      "cmpHash": "default",
      "dagRoundRobin": "disabled",
      "failsafe": "disabled",
      "failsafeAction": "failover-restart-tm",
      "failsafeTimeout": 90,
      "ifIndex": 80,
      "learning": "enable-forward",
      "mtu": 1500,
      "sflow": {
        "pollInterval": 0,
        "pollIntervalGlobal": "yes",
        "samplingRate": 0,
        "samplingRateGlobal": "yes"
      },
      "sourceChecking": "disabled",
      "tag": 10,
      "interfacesReference": {
        "link": "https://localhost/mgmt/tm/net/vlan/~Common~vlan/interfaces?ver=11.5.1",
        "isSubcollection": true
      }
    }]}`))
	}

	vlans, err := s.Client.Vlans()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/vlan", "")
	assert.Equal(s.T(), 1, len(vlans.Vlans))
	assert.Equal(s.T(), "vlan", vlans.Vlans[0].Name)
	assert.Equal(s.T(), 0, vlans.Vlans[0].SFlow.PollInterval)
}

func (s *NetTestSuite) TestCreateVLan() {
	err := s.Client.CreateVlan("name", 1)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/vlan", `{"name":"name", "tag":1, "sflow":{}}`)
}

func (s *NetTestSuite) TestDeleteVLan() {
	err := s.Client.DeleteVlan("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/vlan/name", "")
}

func (s *NetTestSuite) TestModifyVLan() {
	vlan := &Vlan{MTU: 1500}

	err := s.Client.ModifyVlan("name", vlan)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/vlan/name", `{"mtu":1500, "sflow":{}}`)
}

func (s *NetTestSuite) TestRoutes() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:route:routecollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/route?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:route:routestate",
      "name": "default_route",
      "partition": "Common",
      "fullPath": "/Common/default_route",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/route/~Common~default_route?ver=11.5.1",
      "gw": "0.0.0.0",
      "mtu": 0,
      "network": "default"
    }
  ]
}`))
	}

	routes, err := s.Client.Routes()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/route", "")
	assert.Equal(s.T(), 1, len(routes.Routes))
	assert.Equal(s.T(), "default_route", routes.Routes[0].Name)
}

func (s *NetTestSuite) TestGetRoute() {
	resp := `{
		"name": "default_route",
		"partition": "Common",
		"fullPath": "/Common/default_route",
		"gw": "0.0.0.0",
		"mtu": 1500,
		"network": "default"
	}`
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	route, err := s.Client.GetRoute("/Common/default_route")

	s.Require().NoError(err)
	assertRestCall(s, "GET", "/mgmt/tm/net/route/~Common~default_route", "")
	assert.Equal(s.T(), "default_route", route.Name)
	s.requireReserializesTo(resp, route, "Route should reserialize to itself")
}

func (s *NetTestSuite) TestCreateRoute() {
	err := s.Client.CreateRoute("default_route", "default", "0.0.0.0")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "gw":"0.0.0.0"}`)
}

func (s *NetTestSuite) TestAddRoute() {
	err := s.Client.AddRoute(&Route{
		Name:    "default_route",
		Network: "default",
		Gateway: "0.0.0.0",
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "gw":"0.0.0.0"}`)
}

func (s *NetTestSuite) TestAddRouteBlackhole() {
	err := s.Client.AddRoute(&Route{
		Name:      "default_route",
		Network:   "default",
		Blackhole: true,
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "blackhole":true}`)
}

func (s *NetTestSuite) TestAddRouteInterface() {
	err := s.Client.AddRoute(&Route{
		Name:      "default_route",
		Network:   "default",
		Interface: "/Common/internal",
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "tmInterface":"/Common/internal"}`)
}

func (s *NetTestSuite) TestAddRoutePool() {
	err := s.Client.AddRoute(&Route{
		Name:    "default_route",
		Network: "default",
		Pool:    "/Common/test-pool",
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "pool":"/Common/test-pool"}`)
}

func (s *NetTestSuite) TestDeleteRoute() {
	err := s.Client.DeleteRoute("default_route")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/route/default_route", "")
}

func (s *NetTestSuite) TestModifyRoute() {
	route := &Route{Gateway: "1.1.1.1"}

	err := s.Client.ModifyRoute("default_route", route)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/route/default_route", `{"gw":"1.1.1.1"}`)
}

func (s *NetTestSuite) TestRouteDomains() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:route-domain:route-domaincollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/route-domain?ver=11.5.1",
  "items": [
    {
      "kind": "tm:net:route-domain:route-domainstate",
      "name": "0",
      "partition": "Common",
      "fullPath": "/Common/0",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/route-domain/~Common~0?ver=11.5.1",
      "id": 0,
      "strict": "enabled",
      "vlans": [
        "/Common/http-tunnel",
        "/Common/socks-tunnel"
      ]
    }
  ]
}`))
	}

	routes, err := s.Client.RouteDomains()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/route-domain", "")
	assert.Equal(s.T(), 1, len(routes.RouteDomains))
	assert.Equal(s.T(), 2, len(routes.RouteDomains[0].Vlans))
	assert.Equal(s.T(), "0", routes.RouteDomains[0].Name)
}

func (s *NetTestSuite) TestCreateRouteDomain() {
	err := s.Client.CreateRouteDomain("name", 1, false, "vlan1,vlan2")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route-domain", `{"name":"name", "id":1, "vlans":["vlan1","vlan2"], "strict":"disabled"}`)
}

func (s *NetTestSuite) TestDeleteRouteDomain() {
	err := s.Client.DeleteRouteDomain("name")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/route-domain/name", "")
}

func (s *NetTestSuite) TestModifyRouteDomain() {
	route := &RouteDomain{Name: "name", ID: 1, Strict: "enabled"}

	err := s.Client.ModifyRouteDomain("name", route)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/route-domain/name", `{"name":"name", "id": 1, "strict" : "enabled"}`)
}

func (s *NetTestSuite) TestBGPInstances() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:routing:bgp:bgpcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/routing/bgp?ver=13.1.1.2",
  "items": [
    {
      "kind": "tm:net:routing:bgp:bgpstate",
      "name": "test",
      "partition": "Common",
      "fullPath": "/Common/test",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/routing/bgp/~Common~test?ver=13.1.1.2",
      "localAs": 65001
    }
  ]
}`))
	}

	bgpInstances, err := s.Client.BGPInstances()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/routing/bgp", "")
	assert.Equal(s.T(), 1, len(bgpInstances.BGPInstances))
	assert.Equal(s.T(), "test", bgpInstances.BGPInstances[0].Name)
	assert.Equal(s.T(), 65001, bgpInstances.BGPInstances[0].LocalAS)
}

func (s *NetTestSuite) TestCreateBGPInstance() {
	err := s.Client.CreateBGPInstance("name", 65001)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/routing/bgp", `{"name":"name", "localAs":65001}`)
}

func (s *NetTestSuite) TestAddBGPInstance() {
	err := s.Client.AddBGPInstance(&BGPInstance{
		Name:    "test",
		LocalAS: 65001,
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/routing/bgp", `{"name":"test", "localAs":65001}`)
}

func (s *NetTestSuite) TestGetBGPInstance() {
	resp := `{
      "name": "test",
      "partition": "Common",
      "fullPath": "/Common/test",
      "localAs": 65001
	}`
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	bgpInstance, err := s.Client.GetBGPInstance("/Common/test")

	s.Require().NoError(err)
	assertRestCall(s, "GET", "/mgmt/tm/net/routing/bgp/~Common~test", "")
	assert.Equal(s.T(), "test", bgpInstance.Name)
	s.requireReserializesTo(resp, bgpInstance, "BGPInstance should reserialize to itself")
}

func (s *NetTestSuite) TestDeleteBGPInstance() {
	err := s.Client.DeleteBGPInstance("/Common/test")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/routing/bgp/~Common~test", "")
}

func (s *NetTestSuite) TestModifyBGPInstance() {
	bgpInstance := &BGPInstance{Name: "test", LocalAS: 65001}

	err := s.Client.ModifyBGPInstance("/Common/test", bgpInstance)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/routing/bgp/~Common~test", `{"name":"test", "localAs":65001}`)
}

func (s *NetTestSuite) TestBGPNeighbors() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:routing:bgp:neighbor:neighborcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/routing/bgp/~Common~test/neighbor?ver=13.1.1.2",
  "items": [
    {
      "kind": "tm:net:routing:bgp:neighbor:neighborstate",
      "name": "1.1.1.1",
      "fullPath": "1.1.1.1",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/routing/bgp/~Common~test/neighbor/1.1.1.1?ver=13.1.1.2",
      "remoteAs": 65001
    }
  ]
}`))
	}

	bgpNeighbors, err := s.Client.BGPNeighbors("/Common/test")

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor", "")
	assert.Equal(s.T(), 1, len(bgpNeighbors.BGPNeighbors))
	assert.Equal(s.T(), "1.1.1.1", bgpNeighbors.BGPNeighbors[0].Name)
	assert.Equal(s.T(), 65001, bgpNeighbors.BGPNeighbors[0].RemoteAS)
}

func (s *NetTestSuite) TestCreateBGPNeighbor() {
	err := s.Client.CreateBGPNeighbor("/Common/test", "1.1.1.1", 65001)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor", `{"name":"1.1.1.1", "remoteAs":65001}`)
}

func (s *NetTestSuite) TestAddBGPNeighbor() {
	err := s.Client.AddBGPNeighbor("/Common/test", &BGPNeighbor{
		Name:     "1.1.1.1",
		RemoteAS: 65001,
	})

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor", `{"name":"1.1.1.1", "remoteAs":65001}`)
}

func (s *NetTestSuite) TestGetBGPNeighbor() {
	resp := `{
      "name": "1.1.1.1",
      "fullPath": "1.1.1.1",
      "remoteAs": 65001
	}`
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp))
	}

	bgpNeighbor, err := s.Client.GetBGPNeighbor("/Common/test", "1.1.1.1")

	s.Require().NoError(err)
	assertRestCall(s, "GET", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor/1.1.1.1", "")
	assert.Equal(s.T(), "1.1.1.1", bgpNeighbor.Name)
	s.requireReserializesTo(resp, bgpNeighbor, "BGPNeighbor should reserialize to itself")
}

func (s *NetTestSuite) TestDeleteBGPNeighbor() {
	err := s.Client.DeleteBGPNeighbor("/Common/test", "1.1.1.1")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor/1.1.1.1", "")
}

func (s *NetTestSuite) TestModifyBGPNeighbor() {
	bgpNeighbor := &BGPNeighbor{Name: "1.1.1.1", RemoteAS: 65001}

	err := s.Client.ModifyBGPNeighbor("/Common/test", "1.1.1.1", bgpNeighbor)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/routing/bgp/~Common~test/neighbor/1.1.1.1", `{"name":"1.1.1.1", "remoteAs":65001}`)
}

func assertRestCall(s *NetTestSuite, method, path, body string) {
	assert.Equal(s.T(), method, s.LastRequest.Method)
	assert.Equal(s.T(), path, s.LastRequest.URL.Path)
	if body != "" {
		assert.JSONEq(s.T(), body, s.LastRequestBody)
	}
}
