package bigip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
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

	s.Client = NewSession(s.Server.URL, "", "", "", nil)
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
		"POST",
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

func (s *NetTestSuite) TestCreateRoute() {
	err := s.Client.CreateRoute("default_route", "default", "0.0.0.0")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/route", `{"name":"default_route", "network":"default", "gw":"0.0.0.0"}`)
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

func assertRestCall(s *NetTestSuite, method, path, body string) {
	assert.Equal(s.T(), method, s.LastRequest.Method)
	assert.Equal(s.T(), path, s.LastRequest.URL.Path)
	if body != "" {
		assert.JSONEq(s.T(), body, s.LastRequestBody)
	}
}

func (s *NetTestSuite) TestTunnels() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
  "kind": "tm:net:tunnels:tunnel:tunnelcollectionstate",
  "selfLink": "https://localhost/mgmt/tm/net/tunnels/tunnel?ver=13.1.1.2",
  "items": [
    {
      "kind": "tm:net:tunnels:tunnel:tunnelstate",
      "name": "http-tunnel",
      "partition": "Common",
      "fullPath": "/Common/http-tunnel",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/tunnels/tunnel/~Common~http-tunnel?ver=13.1.1.2",
      "autoLasthop": "default",
      "description": "Tunnel for http-explicit profile",
      "idleTimeout": 300,
      "ifIndex": 912,
      "key": 0,
      "localAddress": "any6",
      "mode": "bidirectional",
      "mtu": 0,
      "profile": "/Common/tcp-forward",
      "profileReference": {
        "link": "https://localhost/mgmt/tm/net/tunnels/tcp-forward/~Common~tcp-forward?ver=13.1.1.2"
      },
      "remoteAddress": "any6",
      "secondaryAddress": "any6",
      "tos": "preserve",
      "transparent": "disabled",
      "usePmtu": "enabled"
    },
    {
      "kind": "tm:net:tunnels:tunnel:tunnelstate",
      "name": "socks-tunnel",
      "partition": "Common",
      "fullPath": "/Common/socks-tunnel",
      "generation": 1,
      "selfLink": "https://localhost/mgmt/tm/net/tunnels/tunnel/~Common~socks-tunnel?ver=13.1.1.2",
      "autoLasthop": "default",
      "description": "Tunnel for socks profile",
      "idleTimeout": 300,
      "ifIndex": 928,
      "key": 0,
      "localAddress": "any6",
      "mode": "bidirectional",
      "mtu": 0,
      "profile": "/Common/tcp-forward",
      "profileReference": {
        "link": "https://localhost/mgmt/tm/net/tunnels/tcp-forward/~Common~tcp-forward?ver=13.1.1.2"
      },
      "remoteAddress": "any6",
      "secondaryAddress": "any6",
      "tos": "preserve",
      "transparent": "disabled",
      "usePmtu": "enabled"
    }
  ]
}`))
	}

	tunnels, err := s.Client.Tunnels()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/tunnels/tunnel", "")
	assert.Equal(s.T(), 2, len(tunnels.Tunnels))
	assert.Equal(s.T(), "http-tunnel", tunnels.Tunnels[0].Name)
	assert.Equal(s.T(), "socks-tunnel", tunnels.Tunnels[1].Name)
}

func (s *NetTestSuite) TestGetTunnel() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "autoLasthop": "default",
    "description": "Tunnel for http-explicit profile",
    "fullPath": "/Common/http-tunnel",
    "generation": 1,
    "idleTimeout": 300,
    "ifIndex": 112,
    "key": 0,
    "kind": "tm:net:tunnels:tunnel:tunnelstate",
    "localAddress": "any6",
    "mode": "bidirectional",
    "mtu": 0,
    "name": "http-tunnel",
    "partition": "Common",
    "profile": "/Common/tcp-forward",
    "profileReference": {
        "link": "https://localhost/mgmt/tm/net/tunnels/tcp-forward/~Common~tcp-forward?ver=14.1.0.3"
    },
    "remoteAddress": "any6",
    "secondaryAddress": "any6",
    "selfLink": "https://localhost/mgmt/tm/net/tunnels/tunnel/~Common~http-tunnel?ver=14.1.0.3",
    "tos": "preserve",
    "transparent": "disabled",
    "usePmtu": "enabled"
}`))
	}

	tunnel, err := s.Client.GetTunnel("http-tunnel")

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/tunnels/tunnel/~Common~http-tunnel", "")
	assert.Equal(s.T(), "http-tunnel", tunnel.Name)
	assert.Equal(s.T(), "/Common/tcp-forward", tunnel.Profile)
}

func (s *NetTestSuite) TestCreateTunnel() {
	testTunnel := Tunnel{
		Name:    "some-foo-tunnel",
		Profile: "/Common/some-foo-profile",
	}
	err := s.Client.CreateTunnel(&testTunnel)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/tunnels/tunnel", `{"name":"some-foo-tunnel", "profile":"/Common/some-foo-profile"}`)
}

func (s *NetTestSuite) TestAddTunnel() {
	someTunnel := Tunnel{
		Name:             "foo-tunnel",
		AppService:       "foo-appservice",
		AutoLasthop:      "foo-lasthop",
		Description:      "foo-desc",
		IdleTimeout:      123,
		IfIndex:          456,
		Key:              789,
		LocalAddress:     "foo-local-address",
		Mode:             "foo-mode",
		Mtu:              1440,
		Partition:        "foo-partition",
		Profile:          "foo-profile",
		RemoteAddress:    "foo-remoteaddr",
		SecondaryAddress: "foo-secondaddr",
		Tos:              "foo-tos",
		TrafficGroup:     "foo-tg",
		Transparent:      "foo-transparent",
		UsePmtu:          "foo-pmtu",
	}
	err := s.Client.AddTunnel(&someTunnel)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/tunnels/tunnel", `{"appService":"foo-appservice", "autoLasthop":"foo-lasthop", "description":"foo-desc", "idleTimeout":123, "ifIndex":456, "key":789, "localAddress":"foo-local-address", "mode":"foo-mode", "mtu":1440, "name":"foo-tunnel", "partition":"foo-partition", "profile":"foo-profile", "remoteAddress":"foo-remoteaddr", "secondaryAddress":"foo-secondaddr", "tos":"foo-tos", "trafficGroup":"foo-tg", "transparent":"foo-transparent", "usePmtu":"foo-pmtu"}`)
}

func (s *NetTestSuite) TestDeleteTunnel() {
	err := s.Client.DeleteTunnel("some-foo-tunnel")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/tunnels/tunnel/some-foo-tunnel", "")
}

func (s *NetTestSuite) TestModifyTunnel() {
	tunnel := &Tunnel{Transparent: "enabled"}

	err := s.Client.ModifyTunnel("some-foo-tunnel", tunnel)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/tunnels/tunnel/some-foo-tunnel", `{"transparent":"enabled"}`)
}

var goodVxlansRespnse = `{
    "items": [
	{
            "defaultsFrom": "/Common/vxlan",
            "defaultsFromReference": {
                "link": "https://localhost/mgmt/tm/net/tunnels/vxlan/~Common~vxlan?ver=13.1.1.2"
            },
            "encapsulationType": "vxlan",
            "floodingType": "multipoint",
            "fullPath": "/Common/vxlan-foo",
            "generation": 1,
            "kind": "tm:net:tunnels:vxlan:vxlanstate",
            "name": "vxlan-foo",
            "partition": "foo",
            "port": 4789,
            "selfLink": "https://localhost/mgmt/tm/net/tunnels/vxlan/~foo~vxlan-foo?ver=13.1.1.2"
        },
        {
            "defaultsFrom": "/Common/vxlan",
            "defaultsFromReference": {
                "link": "https://localhost/mgmt/tm/net/tunnels/vxlan/~Common~vxlan?ver=13.1.1.2"
            },
            "encapsulationType": "vxlan",
            "floodingType": "none",
            "fullPath": "/Common/vxlan-bar",
            "generation": 1,
            "kind": "tm:net:tunnels:vxlan:vxlanstate",
            "name": "vxlan-bar",
            "partition": "bar",
            "port": 4789,
            "selfLink": "https://localhost/mgmt/tm/net/tunnels/vxlan/~bar~vxlan-bar?ver=13.1.1.2"
        }
    ],
    "kind": "tm:net:tunnels:vxlan:vxlancollectionstate",
    "selfLink": "https://localhost/mgmt/tm/net/tunnels/vxlan?ver=13.1.1.2"
}`

var goodVxlanRespnse = `{
            "defaultsFrom": "/Common/vxlan",
            "defaultsFromReference": {
                "link": "https://localhost/mgmt/tm/net/tunnels/vxlan/~Common~vxlan?ver=13.1.1.2"
            },
            "encapsulationType": "vxlan",
            "floodingType": "multipoint",
            "fullPath": "/Common/vxlan-foo",
            "generation": 1,
            "kind": "tm:net:tunnels:vxlan:vxlanstate",
            "name": "vxlan-foo",
            "partition": "foo",
            "port": 4789,
            "selfLink": "https://localhost/mgmt/tm/net/tunnels/vxlan/~foo~vxlan-foo?ver=13.1.1.2"
}`

func (s *NetTestSuite) TestVxlans() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(goodVxlansRespnse))
	}

	vxlans, err := s.Client.Vxlans()

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/tunnels/vxlan", "")
	assert.Equal(s.T(), 2, len(vxlans))
	assert.Equal(s.T(), "vxlan-foo", vxlans[0].Name)
	assert.Equal(s.T(), "vxlan-bar", vxlans[1].Name)
}

func (s *NetTestSuite) TestGetVxlan() {
	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(goodVxlanRespnse))
	}

	vxlan, err := s.Client.GetVxlan("~foo~vxlan-foo")

	assert.Nil(s.T(), err)
	assertRestCall(s, "GET", "/mgmt/tm/net/tunnels/vxlan/~foo~vxlan-foo", "")
	assert.Equal(s.T(), "vxlan-foo", vxlan.Name)
	assert.Equal(s.T(), 4789, vxlan.Port)
}

func (s *NetTestSuite) TestCreateVxlan() {
	err := s.Client.CreateVxlan("some-foo-vxlan")

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/tunnels/vxlan", `{"name":"some-foo-vxlan"}`)
}

func (s *NetTestSuite) TestAddVxlan() {
	someVxlan := Vxlan{
		Name:              "foo-vxlan",
		AppService:        "foo-appservice",
		Description:       "foo-desc",
		DefaultsFrom:      "foo-base-profile",
		EncapsulationType: "foo-encap",
		FloodingType:      "foo-ft",
		Partition:         "foo-partition",
		Port:              123,
	}
	err := s.Client.AddVxlan(&someVxlan)

	assert.Nil(s.T(), err)
	assertRestCall(s, "POST", "/mgmt/tm/net/tunnels/vxlan", `{"appService":"foo-appservice", "defaultsFrom":"foo-base-profile", "description":"foo-desc", "encapsulationType":"foo-encap", "floodingType":"foo-ft", "name":"foo-vxlan", "partition":"foo-partition", "port":123}`)
}

func (s *NetTestSuite) TestDeleteVxlan() {
	err := s.Client.DeleteVxlan("some-foo-vxlan")

	assert.Nil(s.T(), err)
	assertRestCall(s, "DELETE", "/mgmt/tm/net/tunnels/vxlan/some-foo-vxlan", "")
}

func (s *NetTestSuite) TestModifyVxlan() {
	vxlan := &Vxlan{Port: 456}

	err := s.Client.ModifyVxlan("some-foo-vxlan", vxlan)

	assert.Nil(s.T(), err)
	assertRestCall(s, "PUT", "/mgmt/tm/net/tunnels/vxlan/some-foo-vxlan", `{"port":456}`)
}
