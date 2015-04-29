package bigip

import (
	"encoding/json"
	"fmt"
)

type Interfaces struct {
	Interfaces []Interface `json:"items"`
}

type Interface struct {
	Name              string `json:"name,omitempty"`
	FullPath          string `json:"fullPath,omitempty"`
	Generation        int    `json:"generation,omitempty"`
	Bundle            string `json:"bundle,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	FlowControl       string `json:"flowControl,omitempty"`
	ForceGigabitFiber string `json:"forceGigabitFiber,omitempty"`
	IfIndex           int    `json:"ifIndex,omitempty"`
	LLDPAdmin         string `json:"lldpAdmin,omitempty"`
	LLDPTlvmap        int    `json:"lldpTlvmap,omitempty"`
	MACAddress        string `json:"macAddress,omitempty"`
	MediaActive       string `json:"mediaActive,omitempty"`
	MediaFixed        string `json:"mediaFixed,omitempty"`
	MediaMax          string `json:"mediaMax,omitempty"`
	MediaSFP          string `json:"mediaSfp,omitempty"`
	MTU               int    `json:"mtu,omitempty"`
	PreferPort        string `json:"preferPort,omitempty"`
	SFlow             struct {
		PollInterval       int    `json:"pollInterval,omitempty"`
		PollIntervalGlobal string `json:"pollIntervalGlobal,omitempty"`
	} `json:"sflow,omitempty"`
	STP             string `json:"stp,omitempty"`
	STPAutoEdgePort string `json:"stpAutoEdgePort,omitempty"`
	STPEdgePort     string `json:"stpEdgePort,omitempty"`
	STPLinkType     string `json:"stpLinkType,omitempty"`
}

type SelfIPs struct {
	SelfIPs []SelfIP `json:"items"`
}

type SelfIP struct {
	Name                  string `json:"name,omitempty"`
	Partition             string `json:"partition,omitempty"`
	FullPath              string `json:"fullPath,omitempty"`
	Generation            int    `json:"generation,omitempty"`
	Address               string `json:"address,omitempty"`
	Floating              string `json:"floating,omitempty"`
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty"`
	TrafficGroup          string `json:"trafficGroup,omitempty"`
	Unit                  int    `json:"unit,omitempty"`
	Vlan                  string `json:"vlan,omitempty"`
	// AllowService          []string `json:"allowService"`
}

type Trunks struct {
	Trunks []Trunk `json:"items"`
}

type Trunk struct {
	Name               string   `json:"name,omitempty"`
	FullPath           string   `json:"fullPath,omitempty"`
	Generation         int      `json:"generation,omitempty"`
	Bandwidth          int      `json:"bandwidth,omitempty"`
	MemberCount        int      `json:"cfgMbrCount,omitempty"`
	DistributionHash   string   `json:"distributionHash,omitempty"`
	ID                 int      `json:"id,omitempty"`
	LACP               string   `json:"lacp,omitempty"`
	LACPMode           string   `json:"lacpMode,omitempty"`
	LACPTimeout        string   `json:"lacpTimeout,omitempty"`
	LinkSelectPolicy   string   `json:"linkSelectPolicy,omitempty"`
	MACAddress         string   `json:"macAddress,omitempty"`
	STP                string   `json:"stp,omitempty"`
	Type               string   `json:"type,omitempty"`
	WorkingMemberCount int      `json:"workingMbrCount,omitempty"`
	Interfaces         []string `json:"interfaces,omitempty"`
}

type Vlans struct {
	Vlans []Vlan `json:"items"`
}

type Vlan struct {
	Name            string `json:"name,omitempty"`
	Partition       string `json:"partition,omitempty"`
	FullPath        string `json:"fullPath,omitempty"`
	Generation      int    `json:"generation,omitempty"`
	AutoLastHop     string `json:"autoLastHop,omitempty"`
	CMPHash         string `json:"cmpHash,omitempty"`
	DAGRoundRobin   string `json:"dagRoundRobin,omitempty"`
	Failsafe        string `json:"failsafe,omitempty"`
	FailsafeAction  string `json:"failsafeAction,omitempty"`
	FailsafeTimeout int    `json:"failsafeTimeout,omitempty"`
	IfIndex         int    `json:"ifIndex,omitempty"`
	Learning        string `json:"learning,omitempty"`
	MTU             int    `json:"mtu,omitempty"`
	SFlow           struct {
		PollInterval       int    `json:"pollInterval,omitempty"`
		PollIntervalGlobal string `json:"pollIntervalGlobal,omitempty"`
		SamplingRate       int    `json:"samplingRate,omitempty"`
		SamplingRateGlobal string `json:"samplingRateGlobal,omitempty"`
	} `json:"sflow,omitempty"`
	SourceChecking string `json:"sourceChecking,omitempty"`
	Tag            int    `json:"tag,omitempty"`
}

type Routes struct {
	Routes []Route `json:"items"`
}

type Route struct {
	Name       string `json:"name,omitempty"`
	Partition  string `json:"partition,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	Gateway    string `json:"gw,omitempty"`
	MTU        int    `json:"mtu,omitempty"`
	Network    string `json:"network,omitempty"`
}

type RouteDomains struct {
	RouteDomains []RouteDomain `json:"items"`
}

type RouteDomain struct {
	Name       string   `json:"name,omitempty"`
	Partition  string   `json:"partition,omitempty"`
	FullPath   string   `json:"fullPath,omitempty"`
	Generation int      `json:"generation,omitempty"`
	ID         int      `json:"id,omitempty"`
	Strict     string   `json:"strict,omitempty"`
	Vlans      []string `json:"vlans,omitempty"`
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

func (b *BigIP) ModifySelfIP(name string, config *SelfIP) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriSelf, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}
	_, err = b.APICall(req)
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

func (b *BigIP) CreateTrunk(name string, interfaces []string, lacp bool) error {
	config := &Trunk{
		Name:       name,
		Interfaces: interfaces,
	}

	if lacp {
		config.LACP = "enabled"
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriTrunk,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}
	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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
