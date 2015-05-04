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
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty"`
	AllowSNAT              string `json:"allowSnat,omitempty"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight,omitempty"`
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
}

// VirtualAddresses contains a list of all virtual addresses on the BIG-IP system.
type VirtualAddresses struct {
	VirtualAddresses []VirtualAddress `json:"items"`
}

// VirtualAddress contains information about each individual virtual address.
type VirtualAddress struct {
	Name                  string `json:"name"`
	Partition             string `json:"partition,omitempty"`
	FullPath              string `json:"fullPath,omitempty"`
	Generation            int    `json:"generation,omitempty"`
	Address               string `json:"address,omitempty"`
	ARP                   string `json:"arp,omitempty"`
	AutoDelete            string `json:"autoDelete,omitempty"`
	ConnectionLimit       int    `json:"connectionLimit,omitempty"`
	Enabled               string `json:"enabled,omitempty"`
	Floating              string `json:"floating,omitempty"`
	ICMPEcho              string `json:"icmpEcho,omitempty"`
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty"`
	Mask                  string `json:"mask,omitempty"`
	RouteAdvertisement    string `json:"routeAdvertisement,omitempty"`
	ServerScope           string `json:"serverScope,omitempty"`
	TrafficGroup          string `json:"trafficGroup,omitempty"`
	Unit                  int    `json:"unit,omitempty"`
}

var (
	uriNode           = "ltm/node"
	uriPolicy         = "ltm/policy"
	uriPool           = "ltm/pool"
	uriVirtual        = "ltm/virtual"
	uriVirtualAddress = "ltm/virtual-address"
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNode removes a node.
func (b *BigIP) DeleteNode(name string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    fmt.Sprintf("%s/%s", uriNode, name),
	}
	_, err := b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// ModifyNode allows you to change any attribute of a node. Fields that
// can be modified are referenced in the Node struct.
func (b *BigIP) ModifyNode(name string, config *Vlan) error {
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
	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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

// AddPoolMembers adds a node/member to the given pool.
func (b *BigIP) AddPoolMember(pool string, member string) error {
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// DeletePoolMember removes a member from the given pool.
func (b *BigIP) DeletePoolMember(pool string, member string) error {
	req := &APIRequest{
		Method:      "delete",
		URL:         fmt.Sprintf("%s/%s/members/%s", uriPool, pool, member),
		ContentType: "application/json",
	}

	_, err := b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// PoolMemberStatus changes the status of a pool member. <state> can be one of
// "enable", "disable", or "offline".
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// DeletePool removes a pool.
func (b *BigIP) DeletePool(name string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    fmt.Sprintf("%s/%s", uriPool, name),
	}
	_, err := b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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
	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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

// CreateVirtualServer adds a new virtual server to the BIG-IP system.
func (b *BigIP) CreateVirtualServer(name, destination, mask, pool string, port int) error {
	var subnetMask string

	if strings.Contains(mask, ".") {
		subnetMask = mask
	} else {
		subnetMask = cidr[mask]
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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteVirtualServer removes a virtual server.
func (b *BigIP) DeleteVirtualServer(name string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    fmt.Sprintf("%s/%s", uriVirtual, name),
	}
	_, err := b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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
	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
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

	switch state {
	case "enable":
		config.Enabled = "yes"
	case "disable":
		config.Enabled = "no"
	}

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

	_, err = b.APICall(req)
	if err != nil {
		return err
	}

	return nil
}
