package bigip

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Nodes contains a list of every node on the BIG-IP system.
type Nodes struct {
	Nodes []Node `json:"items"`
}

// Node contains information about each individual node. You can use all
// of these fields when modifying a node.
type Node struct {
	Name            string `json:"name,omitempty"`
	Partition       string `json:"partition,omitempty"`
	FullPath        string `json:"fullPath,omitempty"`
	Generation      int    `json:"generation,omitempty"`
	Address         string `json:"address,omitempty"`
	ConnectionLimit int    `json:"connectionLimit,omitempty"`
	DynamicRatio    int    `json:"dynamicRatio,omitempty"`
	Logging         string `json:"logging,omitempty"`
	Monitor         string `json:"monitor,omitempty"`
	RateLimit       string `json:"rateLimit,omitempty"`
	Ratio           int    `json:"ratio,omitempty"`
	Session         string `json:"session,omitempty"`
	State           string `json:"state,omitempty"`
}

// Pools contains a list of pools on the BIG-IP system.
type Pools struct {
	Pools []Pool `json:"items"`
}

// Pool contains information about each pool. You can use all of these
// fields when modifying a pool.
type Pool struct {
	Name                   string
	Partition              string
	FullPath               string
	Generation             int
	AllowNAT               bool
	AllowSNAT              bool
	IgnorePersistedWeight  bool
	IPTOSToClient          string
	IPTOSToServer          string
	LinkQoSToClient        string
	LinkQoSToServer        string
	LoadBalancingMode      string
	MinActiveMembers       int
	MinUpMembers           int
	MinUpMembersAction     string
	MinUpMembersChecking   string
	Monitor                string
	QueueDepthLimit        int
	QueueOnConnectionLimit string
	QueueTimeLimit         int
	ReselectTries          int
	SlowRampTime           int
}

// Pool transfer object so we can mask the boolean data munging
type poolDTO struct {
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty" bool:"yes"`
	AllowSNAT              string `json:"allowSnat,omitempty" bool:"yes"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight,omitempty" bool:"enabled"`
	IPTOSToClient          string `json:"ipTosToClient,omitempty"`
	IPTOSToServer          string `json:"ipTosToServer,omitempty"`
	LinkQoSToClient        string `json:"linkQosToClient,omitempty"`
	LinkQoSToServer        string `json:"linkQosToServer,omitempty"`
	LoadBalancingMode      string `json:"loadBalancingMode,omitempty"`
	MinActiveMembers       int    `json:"minActiveMembers,omitempty"`
	MinUpMembers           int    `json:"minUpMembers,omitempty"`
	MinUpMembersAction     string `json:"minUpMembersAction,omitempty"`
	MinUpMembersChecking   string `json:"minUpMembersChecking,omitempty"`
	Monitor                string `json:"monitor"`
	QueueDepthLimit        int    `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int    `json:"queueTimeLimit,omitempty"`
	ReselectTries          int    `json:"reselectTries,omitempty"`
	SlowRampTime           int    `json:"slowRampTime,omitempty"`
}

func (p *Pool) MarshalJSON() ([]byte, error) {
	var dto poolDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Pool) UnmarshalJSON(b []byte) error {
	var dto poolDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

// poolMember is used only when adding members to a pool.
type poolMember struct {
	Name string `json:"name"`
}

// VirtualServers contains a list of all virtual servers on the BIG-IP system.
type VirtualServers struct {
	VirtualServers []VirtualServer `json:"items"`
}

// VirtualServer contains information about each individual virtual server.
type VirtualServer struct {
	Name                     string `json:"name,omitempty"`
	Partition                string `json:"partition,omitempty"`
	FullPath                 string `json:"fullPath,omitempty"`
	Generation               int    `json:"generation,omitempty"`
	AddressStatus            string `json:"addressStatus,omitempty"`
	AutoLastHop              string `json:"autoLastHop,omitempty"`
	CMPEnabled               string `json:"cmpEnabled,omitempty"`
	ConnectionLimit          int    `json:"connectionLimit,omitempty"`
	Destination              string `json:"destination,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
	GTMScore                 int    `json:"gtmScore,omitempty"`
	IPProtocol               string `json:"ipProtocol,omitempty"`
	Mask                     string `json:"mask,omitempty"`
	Mirror                   string `json:"mirror,omitempty"`
	MobileAppTunnel          string `json:"mobileAppTunnel,omitempty"`
	NAT64                    string `json:"nat64,omitempty"`
	Pool                     string `json:"pool,omitempty"`
	RateLimit                string `json:"rateLimit,omitempty"`
	RateLimitDestinationMask int    `json:"rateLimitDstMask,omitempty"`
	RateLimitMode            string `json:"rateLimitMode,omitempty"`
	RateLimitSourceMask      int    `json:"rateLimitSrcMask,omitempty"`
	Source                   string `json:"source,omitempty"`
	SourceAddressTranslation struct {
		Type string `json:"type,omitempty"`
	} `json:"sourceAddressTranslation,omitempty"`
	SourcePort       string   `json:"sourcePort,omitempty"`
	SYNCookieStatus  string   `json:"synCookieStatus,omitempty"`
	TranslateAddress string   `json:"translateAddress,omitempty"`
	TranslatePort    string   `json:"translatePort,omitempty"`
	VlansDisabled    bool     `json:"vlansDisabled,omitempty"`
	VSIndex          int      `json:"vsIndex,omitempty"`
	Rules            []string `json:"rules,omitempty"`
	Profiles		 []Profile `json:"profiles,omitempty"`
}

// VirtualAddresses contains a list of all virtual addresses on the BIG-IP system.
type VirtualAddresses struct {
	VirtualAddresses []VirtualAddress `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type VirtualAddress struct {
	Name                  string
	Partition             string
	FullPath              string
	Generation            int
	Address               string
	ARP                   bool
	AutoDelete            string
	ConnectionLimit       int
	Enabled               bool
	Floating              bool
	ICMPEcho              bool
	InheritedTrafficGroup bool
	Mask                  string
	RouteAdvertisement    bool
	ServerScope           string
	TrafficGroup          string
	Unit                  int
}

type virtualAddressDTO struct {
	Name                  string `json:"name"`
	Partition             string `json:"partition,omitempty"`
	FullPath              string `json:"fullPath,omitempty"`
	Generation            int    `json:"generation,omitempty"`
	Address               string `json:"address,omitempty"`
	ARP                   string `json:"arp,omitempty" bool:"enabled"`
	AutoDelete            string `json:"autoDelete,omitempty"`
	ConnectionLimit       int    `json:"connectionLimit,omitempty"`
	Enabled               string `json:"enabled,omitempty" bool:"yes"`
	Floating              string `json:"floating,omitempty" bool:"enabled"`
	ICMPEcho              string `json:"icmpEcho,omitempty" bool:"enabled"`
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty" bool:"yes"`
	Mask                  string `json:"mask,omitempty"`
	RouteAdvertisement    string `json:"routeAdvertisement,omitempty" bool:"enabled"`
	ServerScope           string `json:"serverScope,omitempty"`
	TrafficGroup          string `json:"trafficGroup,omitempty"`
	Unit                  int    `json:"unit,omitempty"`
}

func (p *VirtualAddress) MarshalJSON() ([]byte, error) {
	var dto virtualAddressDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *VirtualAddress) UnmarshalJSON(b []byte) error {
	var dto virtualAddressDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

// Monitors contains a list of all monitors on the BIG-IP system.
type Monitors struct {
	Monitors []Monitor `json:"items"`
}

// Monitor contains information about each individual monitor.
type Monitor struct {
	Name           string
	Partition      string
	FullPath       string
	Generation     int
	ParentMonitor  string
	Description    string
	Destination    string
	Interval       int
	IPDSCP         int
	ManualResume   bool
	ReceiveString  string
	ReceiveDisable string
	Reverse        bool
	SendString     string
	TimeUntilUp    int
	Timeout        int
	Transparent    bool
	UpInterval     int
}

type monitorDTO struct {
	Name           string `json:"name,omitempty"`
	Partition      string `json:"partition,omitempty"`
	FullPath       string `json:"fullPath,omitempty"`
	Generation     int    `json:"generation,omitempty"`
	ParentMonitor  string `json:"defaultsFrom,omitempty"`
	Description    string `json:"description,omitempty"`
	Destination    string `json:"destination,omitempty"`
	Interval       int    `json:"interval,omitempty"`
	IPDSCP         int    `json:"ipDscp,omitempty"`
	ManualResume   string `json:"manualResume,omitempty" bool:"enabled"`
	ReceiveString  string `json:"recv,omitempty"`
	ReceiveDisable string `json:"recvDisable,omitempty"`
	Reverse        string `json:"reverse,omitempty" bool:"enabled"`
	SendString     string `json:"send,omitempty"`
	TimeUntilUp    int    `json:"timeUntilUp,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Transparent    string `json:"transparent,omitempty" bool:"enabled"`
	UpInterval     int    `json:"upInterval,omitempty"`
}

type Profiles struct {
	Profiles []Profile `json:"items"`
}

type Profile struct {
	Name		   string `json:"name,omitempty"`
	Partition      string `json:"partition,omitempty"`
}

type IRules struct {
	IRules	[]IRule `json:"items"`
}

type IRule struct {
	Name		   string `json:"name,omitempty"`
	Partition      string `json:"partition,omitempty"`
	Rule		   string `json:"apiAnonymous,omitempty"`
}

func (p *Monitor) MarshalJSON() ([]byte, error) {
	var dto monitorDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *Monitor) UnmarshalJSON(b []byte) error {
	var dto monitorDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}

var (
	uriNode           = "ltm/node"
	uriPolicy         = "ltm/policy"
	uriPool           = "ltm/pool"
	uriVirtual        = "ltm/virtual"
	uriVirtualAddress = "ltm/virtual-address"
	uriMonitor        = "ltm/monitor"
	uriRule			  = "ltm/rule"
	cidr              = map[string]string{
		"0":  "0.0.0.0",
		"1":  "128.0.0.0",
		"2":  "192.0.0.0",
		"3":  "224.0.0.0",
		"4":  "240.0.0.0",
		"5":  "248.0.0.0",
		"6":  "252.0.0.0",
		"7":  "254.0.0.0",
		"8":  "255.0.0.0",
		"9":  "255.128.0.0",
		"10": "255.192.0.0",
		"11": "255.224.0.0",
		"12": "255.240.0.0",
		"13": "255.248.0.0",
		"14": "255.252.0.0",
		"15": "255.254.0.0",
		"16": "255.255.0.0",
		"17": "255.255.128.0",
		"18": "255.255.192.0",
		"19": "255.255.224.0",
		"20": "255.255.240.0",
		"21": "255.255.248.0",
		"22": "255.255.252.0",
		"23": "255.255.254.0",
		"24": "255.255.255.0",
		"25": "255.255.255.128",
		"26": "255.255.255.192",
		"27": "255.255.255.224",
		"28": "255.255.255.240",
		"29": "255.255.255.248",
		"30": "255.255.255.252",
		"31": "255.255.255.254",
		"32": "255.255.255.255",
	}
)

// Nodes returns a list of nodes.
func (b *BigIP) Nodes() (*Nodes, error) {
	var nodes Nodes
	req := &APIRequest{
		Method: "get",
		URL:    uriNode,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &nodes)
	if err != nil {
		return nil, err
	}

	return &nodes, nil
}

// CreateNode adds a new node to the BIG-IP system.
func (b *BigIP) CreateNode(name, address string) error {
	config := &Node{
		Name:    name,
		Address: address,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriNode,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Get a Node by name. Returns nil if the node does not exist
func (b *BigIP) GetNode(name string) (*Node, error) {
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s", uriNode, name))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var node Node
	err = json.Unmarshal(resp, &node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

// DeleteNode removes a node.
func (b *BigIP) DeleteNode(name string) error {
	return b.delete(uriNode, name)
}

// ModifyNode allows you to change any attribute of a node. Fields that
// can be modified are referenced in the Node struct.
func (b *BigIP) ModifyNode(name string, config *Node) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriNode, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// NodeStatus changes the status of a node. <state> can be either
// "enable" or "disable".
func (b *BigIP) NodeStatus(name, state string) error {
	config := &Node{}

	switch state {
	case "enable":
		// config.State = "unchecked"
		config.Session = "user-enabled"
	case "disable":
		// config.State = "unchecked"
		config.Session = "user-disabled"
		// case "offline":
		// 	config.State = "user-down"
		// 	config.Session = "user-disabled"
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriNode, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Pools returns a list of pools.
func (b *BigIP) Pools() (*Pools, error) {
	var pools Pools
	req := &APIRequest{
		Method: "get",
		URL:    uriPool,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &pools)
	if err != nil {
		return nil, err
	}

	return &pools, nil
}

// PoolMembers returns a list of pool members for the given pool.
func (b *BigIP) PoolMembers(name string) ([]string, error) {
	var nodes Nodes
	members := []string{}
	errString := []string{}
	req := &APIRequest{
		Method: "get",
		URL:    fmt.Sprintf("%s/%s/members", uriPool, name),
	}

	resp, err := b.APICall(req)
	if err != nil {
		return errString, err
	}

	err = json.Unmarshal(resp, &nodes)
	if err != nil {
		return errString, err
	}

	for _, m := range nodes.Nodes {
		members = append(members, m.Name)
	}

	return members, nil
}

// AddPoolMember adds a node/member to the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) AddPoolMember(pool, member string) error {
	config := &poolMember{
		Name: member,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         fmt.Sprintf("%s/%s/members", uriPool, pool),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// DeletePoolMember removes a member from the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) DeletePoolMember(pool, member string) error {
	return b.delete(uriPool, pool, "members", member)
}

// PoolMemberStatus changes the status of a pool member. <state> can be either
// "enable" or "disable". <member> must be in the form of <node>:<port>,
// i.e.: "web-server1:443".
func (b *BigIP) PoolMemberStatus(pool, member, state string) error {
	config := &Node{}

	switch state {
	case "enable":
		// config.State = "unchecked"
		config.Session = "user-enabled"
	case "disable":
		// config.State = "unchecked"
		config.Session = "user-disabled"
		// case "offline":
		// 	config.State = "user-down"
		// 	config.Session = "user-disabled"
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s/members/%s", uriPool, pool, member),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// CreatePool adds a new pool to the BIG-IP system.
func (b *BigIP) CreatePool(name string) error {
	config := &Pool{
		Name: name,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriPool,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Get a Pool by name. Returns nil if the Pool does not exist
func (b *BigIP) GetPool(name string) (*Pool, error) {
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s", uriPool, name))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var pool Pool
	err = json.Unmarshal(resp, &pool)
	if err != nil {
		return nil, err
	}

	return &pool, nil
}

// DeletePool removes a pool.
func (b *BigIP) DeletePool(name string) error {
	return b.delete(uriPool, name)
}

// ModifyPool allows you to change any attribute of a pool. Fields that
// can be modified are referenced in the Pool struct.
func (b *BigIP) ModifyPool(name string, config *Pool) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriPool, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// VirtualServers returns a list of virtual servers.
func (b *BigIP) VirtualServers() (*VirtualServers, error) {
	var vs VirtualServers
	req := &APIRequest{
		Method: "get",
		URL:    uriVirtual,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &vs)
	if err != nil {
		return nil, err
	}

	return &vs, nil
}

// CreateVirtualServer adds a new virtual server to the BIG-IP system. <mask> can either be
// in CIDR notation or decimal, i.e.: "24" or "255.255.255.0". A CIDR mask of "0" is the same
// as "0.0.0.0".
func (b *BigIP) CreateVirtualServer(name, destination, mask, pool string, port int) error {
	subnetMask := cidr[mask]

	if strings.Contains(mask, ".") {
		subnetMask = mask
	}

	config := &VirtualServer{
		Name:        name,
		Destination: fmt.Sprintf("%s:%d", destination, port),
		Mask:        subnetMask,
		Pool:        pool,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uriVirtual,
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Get a VirtualServer by name. Returns nil if the VirtualServer does not exist
func (b *BigIP) GetVirtualServer(name string) (*VirtualServer, error) {
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s", uriVirtual, name))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var vs VirtualServer
	err = json.Unmarshal(resp, &vs)
	if err != nil {
		return nil, err
	}

	profiles,err := b.VirtualServerProfiles(name)
	if err != nil {
		return nil, err
	}

	vs.Profiles = profiles.Profiles
	return &vs, nil
}

// DeleteVirtualServer removes a virtual server.
func (b *BigIP) DeleteVirtualServer(name string) error {
	return b.delete(uriVirtual, name)
}

// ModifyVirtualServer allows you to change any attribute of a virtual server. Fields that
// can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyVirtualServer(name string, config *VirtualServer) error {
	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriVirtual, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// VirtualServerProfiles gets the profiles currently associated with a virtual server.
func (b *BigIP) VirtualServerProfiles(vs string) (*Profiles, error){
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s/profiles", uriVirtual, vs))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var p Profiles
	err = json.Unmarshal(resp, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// VirtualAddresses returns a list of virtual addresses.
func (b *BigIP) VirtualAddresses() (*VirtualAddresses, error) {
	var va VirtualAddresses
	req := &APIRequest{
		Method: "get",
		URL:    uriVirtualAddress,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &va)
	if err != nil {
		return nil, err
	}

	return &va, nil
}

// VirtualAddressStatus changes the status of a virtual address. <state> can be either
// "enable" or "disable".
func (b *BigIP) VirtualAddressStatus(vaddr, state string) error {
	config := &VirtualAddress{}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriVirtualAddress, vaddr),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// Monitors returns a list of all HTTP, HTTPS, Gateway ICMP, and ICMP monitors.
func (b *BigIP) Monitors() ([]Monitor, error) {
	var monitors []Monitor
	monitorUris := []string{"http", "https", "icmp", "gateway-icmp"}

	for _, name := range monitorUris {
		var m Monitors
		req := &APIRequest{
			Method: "get",
			URL:    fmt.Sprintf("%s/%s", uriMonitor, name),
		}

		resp, err := b.APICall(req)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp, &m)
		if err != nil {
			return nil, err
		}

		for _, monitor := range m.Monitors {
			monitors = append(monitors, monitor)
		}
	}

	return monitors, nil
}

// CreateMonitor adds a new monitor to the BIG-IP system. <parent> must be one of "http", "https",
// "icmp", or "gateway icmp".
func (b *BigIP) CreateMonitor(name, parent string, interval, timeout int, send, receive string) error {
	if strings.Contains(send, "\r\n") {
		send = strings.Replace(send, "\r\n", "\\r\\n", -1)
	}

	if strings.Contains(parent, "gateway") {
		parent = "gateway_icmp"
	}

	config := &Monitor{
		Name:          name,
		ParentMonitor: parent,
		Interval:      interval,
		Timeout:       timeout,
		SendString:    send,
		ReceiveString: receive,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         fmt.Sprintf("%s/%s", uriMonitor, parent),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// DeleteMonitor removes a monitor.
func (b *BigIP) DeleteMonitor(name, parent string) error {
	return b.delete(uriMonitor, parent, name)
}

// ModifyMonitor allows you to change any attribute of a monitor. <parent> must be
// one of "http", "https", "icmp", or "gateway icmp". Fields that
// can be modified are referenced in the Monitor struct.
func (b *BigIP) ModifyMonitor(name, parent string, config *Monitor) error {
	if strings.Contains(config.SendString, "\r\n") {
		config.SendString = strings.Replace(config.SendString, "\r\n", "\\r\\n", -1)
	}

	if strings.Contains(parent, "gateway") {
		parent = "gateway_icmp"
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s/%s", uriMonitor, parent, name),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// AddMonitorToPool assigns the monitor, <monitor> to the given <pool>.
func (b *BigIP) AddMonitorToPool(monitor, pool string) error {
	config := &Pool{
		Monitor: monitor,
	}

	marshalJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("%s/%s", uriPool, pool),
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

// IRules returns a list of irules
func (b *BigIP) IRules() (*IRules, error) {
	var rules IRules
	req := &APIRequest{
		Method: "get",
		URL:    uriRule,
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &rules)
	if err != nil {
		return nil, err
	}

	return &rules, nil
}

func (b *BigIP) IRule(name string) (*IRule, error) {
	var rule IRule
	req := &APIRequest{
		Method: "get",
		URL:    fmt.Sprintf("%s/%s", uriRule, name),
	}

	resp, err := b.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &rule)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// IRules returns a list of irules
func (b *BigIP) CreateIRule(name, rule string) error {
	irule :=  &IRule{
		Name: name,
		Rule: rule,
	}
	return b.post(irule, uriRule)
}

// DeleteNode removes a node.
func (b *BigIP) DeleteIRule(name string) error {
	return b.delete(uriRule, name)
}

func (b *BigIP) ModifyIRule(name string, irule *IRule) error {
	irule.Name = name
	return b.put(irule, uriRule, name)
}