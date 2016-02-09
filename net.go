package bigip

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Interfaces contains a list of every interface on the BIG-IP system.
type Interfaces struct {
	Interfaces []Interface `json:"items"`
}

// Interface contains information about each individual interface.
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

// SelfIPs contains a list of every self IP on the BIG-IP system.
type SelfIPs struct {
	SelfIPs []SelfIP `json:"items"`
}

// SelfIP contains information about each individual self IP. You can use all of
// these fields when modifying a self IP.
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

// Trunks contains a list of every trunk on the BIG-IP system.
type Trunks struct {
	Trunks []Trunk `json:"items"`
}

// Trunk contains information about each individual trunk. You can use all of
// these fields when modifying a trunk.
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

// Vlans contains a list of every VLAN on the BIG-IP system.
type Vlans struct {
	Vlans []Vlan `json:"items"`
}

// Vlan contains information about each individual VLAN. You can use all of
// these fields when modifying a VLAN.
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

// VlanInterface contains fields to be used when adding an interface to a VLAN.
type VlanInterface struct {
	Name     string `json:"name,omitempty"`
	Tagged   bool   `json:"tagged,omitempty"`
	Untagged bool   `json:"untagged,omitempty"`
}

// Routes contains a list of every route on the BIG-IP system.
type Routes struct {
	Routes []Route `json:"items"`
}

// Route contains information about each individual route. You can use all
// of these fields when modifying a route.
type Route struct {
	Name       string `json:"name,omitempty"`
	Partition  string `json:"partition,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	Gateway    string `json:"gw,omitempty"`
	MTU        int    `json:"mtu,omitempty"`
	Network    string `json:"network,omitempty"`
}

// RouteDomains contains a list of every route domain on the BIG-IP system.
type RouteDomains struct {
	RouteDomains []RouteDomain `json:"items"`
}

// RouteDomain contains information about each individual route domain. You can use all
// of these fields when modifying a route domain.
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

// Interfaces returns a list of interfaces.
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

// AddInterfaceToVlan associates the given interface to the specified VLAN.
func (b *BigIP) AddInterfaceToVlan(vlan, iface string, tagged bool) error {
	config := &VlanInterface{}

	if tagged {
		config.Name = vlan
		config.Tagged = true
	}

	if !tagged {
		config.Name = vlan
		config.Untagged = true
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s/interfaces", uriVlan, vlan),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// SelfIPs returns a list of self IP's.
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

// CreateSelfIP adds a new self IP to the BIG-IP system. For <address>, you
// must include the subnet mask in CIDR notation, i.e.: "10.1.1.1/24".
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

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteSelfIP removes a self IP.
func (b *BigIP) DeleteSelfIP(name string) error {
	return b.delete(uriSelf, name)
}

// ModifySelfIP allows you to change any attribute of a self IP. Fields that
// can be modified are referenced in the SelfIP struct.
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

	_, callErr := b.APICall(req)
	return callErr
}

// Trunks returns a list of trunks.
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

// CreateTrunk adds a new trunk to the BIG-IP system. <interfaces> must be
// separated by a comma, i.e.: "1.4, 1.6, 1.8".
func (b *BigIP) CreateTrunk(name, interfaces string, lacp bool) error {
	rawInts := strings.Split(interfaces, ",")
	ints := []string{}

	for _, i := range rawInts {
		ints = append(ints, strings.Trim(i, " "))
	}

	config := &Trunk{
		Name:       name,
		Interfaces: ints,
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

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteTrunk removes a trunk.
func (b *BigIP) DeleteTrunk(name string) error {
	return b.delete(uriTrunk, name)
}

// ModifyTrunk allows you to change any attribute of a trunk. Fields that
// can be modified are referenced in the Trunk struct.
func (b *BigIP) ModifyTrunk(name string, config *Trunk) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriTrunk, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Vlans returns a list of vlans.
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

// CreateVlan adds a new VLAN to the BIG-IP system.
func (b *BigIP) CreateVlan(name string, tag int) error {
	config := &Vlan{
		Name: name,
		Tag:  tag,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriVlan,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteVlan removes a vlan.
func (b *BigIP) DeleteVlan(name string) error {
	return b.delete(uriVlan, name)
}

// ModifyVlan allows you to change any attribute of a VLAN. Fields that
// can be modified are referenced in the Vlan struct.
func (b *BigIP) ModifyVlan(name string, config *Vlan) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriVlan, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Routes returns a list of routes.
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

// CreateRoute adds a new static route to the BIG-IP system. <dest> must include the
// subnet mask in CIDR notation, i.e.: "10.1.1.0/24".
func (b *BigIP) CreateRoute(name, dest, gateway string) error {
	config := &Route{
		Name:    name,
		Network: dest,
		Gateway: gateway,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriRoute,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteRoute removes a static route.
func (b *BigIP) DeleteRoute(name string) error {
	return b.delete(uriRoute, name)
}

// ModifyRoute allows you to change any attribute of a static route. Fields that
// can be modified are referenced in the Route struct.
func (b *BigIP) ModifyRoute(name string, config *Route) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriRoute, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// RouteDomains returns a list of route domains.
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

// CreateRouteDomain adds a new route domain to the BIG-IP system. <vlans> must be separated
// by a comma, i.e.: "vlan1010, vlan1020".
func (b *BigIP) CreateRouteDomain(name string, id int, strict bool, vlans string) error {
	strictIsolation := "enabled"
	vlanMembers := []string{}
	rawVlans := strings.Split(vlans, ",")

	for _, v := range rawVlans {
		vlanMembers = append(vlanMembers, strings.Trim(v, " "))
	}

	if !strict {
		strictIsolation = "disabled"
	}

	config := &RouteDomain{
		Name:   name,
		ID:     id,
		Strict: strictIsolation,
		Vlans:  vlanMembers,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriRouteDomain,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteRouteDomain removes a route domain.
func (b *BigIP) DeleteRouteDomain(name string) error {
	return b.delete(uriRouteDomain, name)
}

// ModifyRouteDomain allows you to change any attribute of a route domain. Fields that
// can be modified are referenced in the RouteDomain struct.
func (b *BigIP) ModifyRouteDomain(name string, config *RouteDomain) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriRouteDomain, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}
