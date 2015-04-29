package bigip

import (
	"encoding/json"
	"fmt"
)

type Interfaces struct {
	Interfaces []Interface `json:"items"`
}

type Interface struct {
	Name              string `json:"name"`
	FullPath          string `json:"fullPath"`
	Generation        int    `json:"generation"`
	Bundle            string `json:"bundle"`
	Enabled           bool   `json:"enabled"`
	FlowControl       string `json:"flowControl"`
	ForceGigabitFiber string `json:"forceGigabitFiber"`
	IfIndex           int    `json:"ifIndex"`
	LLDPAdmin         string `json:"lldpAdmin"`
	LLDPTlvmap        int    `json:"lldpTlvmap"`
	MACAddress        string `json:"macAddress"`
	MediaActive       string `json:"mediaActive"`
	MediaFixed        string `json:"mediaFixed"`
	MediaMax          string `json:"mediaMax"`
	MediaSFP          string `json:"mediaSfp"`
	MTU               int    `json:"mtu"`
	PreferPort        string `json:"preferPort"`
	SFlow             struct {
		PollInterval       int    `json:"pollInterval"`
		PollIntervalGlobal string `json:"pollIntervalGlobal"`
	} `json:"sflow"`
	STP             string `json:"stp"`
	STPAutoEdgePort string `json:"stpAutoEdgePort"`
	STPEdgePort     string `json:"stpEdgePort"`
	STPLinkType     string `json:"stpLinkType"`
}

type SelfIPs struct {
	SelfIPs []SelfIP `json:"items"`
}

type SelfIP struct {
	Name                  string `json:"name"`
	Partition             string `json:"partition"`
	FullPath              string `json:"fullPath"`
	Generation            int    `json:"generation"`
	Address               string `json:"address"`
	Floating              string `json:"floating"`
	InheritedTrafficGroup string `json:"inheritedTrafficGroup"`
	TrafficGroup          string `json:"trafficGroup"`
	Unit                  int    `json:"unit"`
	Vlan                  string `json:"vlan"`
	// AllowService          []string `json:"allowService"`
}

type Trunks struct {
	Trunks []Trunk `json:"items"`
}

type Trunk struct {
	Name               string   `json:"name"`
	FullPath           string   `json:"fullPath"`
	Generation         int      `json:"generation"`
	Bandwidth          int      `json:"bandwidth"`
	MemberCount        int      `json:"cfgMbrCount"`
	DistributionHash   string   `json:"distributionHash"`
	ID                 int      `json:"id"`
	LACP               string   `json:"lacp"`
	LACPMode           string   `json:"lacpMode"`
	LACPTimeout        string   `json:"lacpTimeout"`
	LinkSelectPolicy   string   `json:"linkSelectPolicy"`
	MACAddress         string   `json:"macAddress"`
	STP                string   `json:"stp"`
	Type               string   `json:"type"`
	WorkingMemberCount int      `json:"workingMbrCount"`
	Interfaces         []string `json:"interfaces"`
}

type Vlans struct {
	Vlans []Vlan `json:"items"`
}

type Vlan struct {
	Name            string `json:"name"`
	Partition       string `json:"partition"`
	FullPath        string `json:"fullPath"`
	Generation      int    `json:"generation"`
	AutoLastHop     string `json:"autoLastHop"`
	CMPHash         string `json:"cmpHash"`
	DAGRoundRobin   string `json:"dagRoundRobin"`
	Failsafe        string `json:"failsafe"`
	FailsafeAction  string `json:"failsafeAction"`
	FailsafeTimeout int    `json:"failsafeTimeout"`
	IfIndex         int    `json:"ifIndex"`
	Learning        string `json:"learning"`
	MTU             int    `json:"mtu"`
	SFlow           struct {
		PollInterval       int    `json:"pollInterval"`
		PollIntervalGlobal string `json:"pollIntervalGlobal"`
		SamplingRate       int    `json:"samplingRate"`
		SamplingRateGlobal string `json:"samplingRateGlobal"`
	} `json:"sflow"`
	SourceChecking string `json:"sourceChecking"`
	Tag            int    `json:"tag"`
}

type Routes struct {
	Routes []Route `json:"items"`
}

type Route struct {
	Name       string `json:"name"`
	Partition  string `json:"partition"`
	FullPath   string `json:"fullPath"`
	Generation int    `json:"generation"`
	Gateway    string `json:"gw"`
	MTU        int    `json:"mtu"`
	Network    string `json:"network"`
}

type RouteDomains struct {
	RouteDomains []RouteDomain `json:"items"`
}

type RouteDomain struct {
	Name       string   `json:"name"`
	Partition  string   `json:"partition"`
	FullPath   string   `json:"fullPath"`
	Generation int      `json:"generation"`
	ID         int      `json:"id"`
	Strict     string   `json:"strict"`
	Vlans      []string `json:"vlans"`
}

var (
	uriInterface   = "net/interface"
	uriSelf        = "net/self"
	uriTrunk       = "net/trunk"
	uriVlan        = "net/vlan"
	uriRoute       = "net/route"
	uriRouteDomain = "net/route-domain"
)

func (b *BigIP) Interfaces() (*Interfaces, error) {
	var interfaces Interfaces
	req := &APIRequest{
		Method: "get",
		URL:    uriInterface,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &interfaces)
	if err != nil {
		return nil, err
	}

	return &interfaces, nil
}

func (b *BigIP) SelfIPs() (*SelfIPs, error) {
	var self SelfIPs
	req := &APIRequest{
		Method: "get",
		URL:    uriSelf,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &self)
	if err != nil {
		return nil, err
	}

	return &self, nil
}

func (b *BigIP) CreateSelfIP(name, address, vlan string) error {
	config := &SelfIP{
		Name:    name,
		Address: address,
		Vlan:    vlan,
	}
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriSelf,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}
	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

func (b *BigIP) DeleteSelfIP(name string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    fmt.Sprintf("%s/%s", uriSelf, name),
	}
	_, err := b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

func (b *BigIP) Trunks() (*Trunks, error) {
	var trunks Trunks
	req := &APIRequest{
		Method: "get",
		URL:    uriTrunk,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &trunks)
	if err != nil {
		return nil, err
	}

	return &trunks, nil
}

func (b *BigIP) Vlans() (*Vlans, error) {
	var vlans Vlans
	req := &APIRequest{
		Method: "get",
		URL:    uriVlan,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &vlans)
	if err != nil {
		return nil, err
	}

	return &vlans, nil
}

func (b *BigIP) Routes() (*Routes, error) {
	var routes Routes
	req := &APIRequest{
		Method: "get",
		URL:    uriRoute,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &routes)
	if err != nil {
		return nil, err
	}

	return &routes, nil
}

func (b *BigIP) RouteDomains() (*RouteDomains, error) {
	var rd RouteDomains
	req := &APIRequest{
		Method: "get",
		URL:    uriRouteDomain,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &rd)
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
