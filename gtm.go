package bigip

import (
	"fmt"
)

// GTM Documentation
// https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_gtm.ashx

// ********************************************************************************************************************
// *************************************************                  *************************************************
// *************************************************   GTM WideIP A   *************************************************
// *************************************************                  *************************************************
// ********************************************************************************************************************

// GTMWideIPs contains a list of every WideIP on the BIG-IP system.
type GTMWideIPs struct {
	GTMWideIPs []GTMWideIP `json:"items"`
}

// GTMWideIP contains information about each wide ip (regardless of type: A, AAAA, CNAME, etc.)
// Type is what determine the type of record the WideIp is for in the docs, however that is NOT returned by the API
// Instead you have to query the Type by the uri   wideip/a  wideip/cname  that = type
type GTMWideIP struct {
	Name                              string   `json:"name,omitempty"`
	Partition                         string   `json:"partition,omitempty"`
	FullPath                          string   `json:"fullPath,omitempty"`
	Generation                        int      `json:"generation,omitempty"`
	AppService                        string   `json:"appService,omitempty"`
	Description                       string   `json:"description,omitempty"`
	Disabled                          bool     `json:"disabled,omitempty"`
	Enabled                           bool     `json:"enabled,omitempty"`
	FailureRcode                      string   `json:"failureRcode,omitempty"`
	FailureRcodeResponse              string   `json:"failureRcodeResponse,omitempty"`
	FailureRcodeTTL                   int      `json:"failureRcodeTtl,omitempty"`
	LastResortPool                    string   `json:"lastResortPool,omitempty"`
	LoadBalancingDecisionLogVerbosity []string `json:"loadBalancingDecisionLogVerbosity,omitempty"`
	MinimalResponse                   string   `json:"minimalResponse,omitempty"`
	PersistCidrIpv4                   int      `json:"persistCidrIpv4,omitempty"`
	PersistCidrIpv6                   int      `json:"persistCidrIpv6,omitempty"`
	Persistence                       string   `json:"persistence,omitempty"`
	PoolLbMode                        string   `json:"poolLbMode,omitempty"`
	TTLPersistence                    int      `json:"ttlPersistence,omitempty"`

	// Not in the spec, but returned by the API
	// Setting this field atomically updates all members.
	Pools *[]GTMWideIPPool `json:"pools,omitempty"`
}

// GTMWideIPPool Pool Structure
type GTMWideIPPool struct {
	Name          string `json:"name,omitempty"`
	Partition     string `json:"partition,omitempty"`
	Order         int    `json:"order,omitempty"`
	Ratio         int    `json:"ratio,omitempty"`
	NameReference struct {
		Link string `json:"link,omitempty"`
	} `json:"nameReference,omitempty"`
}

// GetGTMWideIPs returns a list of all WideIps for a provided type
func (b *BigIP) GetGTMWideIPs(recordType GTMType) (*GTMWideIPs, error) {
	var w GTMWideIPs
	err, _ := b.getForEntity(&w, uriGtm, uriWideIp, string(recordType))
	if err != nil {
		return nil, err
	}

	return &w, nil
}

// GetGTMWideIP get's a WideIP by name
func (b *BigIP) GetGTMWideIP(name string, recordType GTMType) (*GTMWideIP, error) {
	var w GTMWideIP

	err, ok := b.getForEntity(&w, uriGtm, uriWideIp, string(recordType), name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}

// AddGTMWideIP adds a WideIp by config to the BIG-IP system.
func (b *BigIP) AddGTMWideIP(config *GTMWideIP, recordType GTMType) error {
	return b.post(config, uriGtm, uriWideIp, string(recordType))
}

// DeleteGTMWideIP removes a WideIp by config to the BIG-IP system.
func (b *BigIP) DeleteGTMWideIP(fullPath string, recordType GTMType) error {
	return b.delete(uriGtm, uriWideIp, string(recordType), fullPath)
}

// ModifyGTMWideIP adds a WideIp by config to the BIG-IP system.
func (b *BigIP) ModifyGTMWideIP(fullPath string, config *GTMWideIP, recordType GTMType) error {
	return b.put(config, uriGtm, uriWideIp, string(recordType), fullPath)
}

// ********************************************************************************************************************
// ********************************************                     ***************************************************
// ********************************************   GTM Pool Common   ***************************************************
// ********************************************                     ***************************************************
// ********************************************************************************************************************

// DeleteGTMPool removes a Pool by config and Pool Type from the BIG-IP system.
func (b *BigIP) DeleteGTMPool(fullPath string, recordType GTMType) error {
	return b.delete(uriGtm, uriPool, string(recordType), fullPath)
}

// ********************************************************************************************************************
// **********************************************                ******************************************************
// **********************************************   GTM Pool A   ******************************************************
// **********************************************                ******************************************************
// ********************************************************************************************************************

// GTMAPools contains a list of every gtm/pool/a on the BIG-IP system.
type GTMAPools struct {
	GTMAPools []GTMAPool `json:"items"`
}

// GTMAPool contains information about each gtm/pool/a
type GTMAPool struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AppService                string `json:"appService,omitempty"`
	Description               string `json:"description,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	DynamicRatio              string `json:"dynamicRatio,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	FallbackIP                string `json:"fallbackIp,omitempty"`
	FallbackMode              string `json:"fallbackMode,omitempty"`
	LimitMaxBps               uint64 `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       uint64 `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               uint64 `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	LoadBalancingMode         string `json:"loadBalancingMode,omitempty"`
	ManualResume              string `json:"manualResume,omitempty"`
	MaxAnswersReturned        int    `json:"maxAnswersReturned,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	TmPartition               string `json:"tmPartition,omitempty"`
	QosHitRatio               int    `json:"qosHitRatio,omitempty"`
	QosHops                   int    `json:"qosHops,omitempty"`
	QosKilobytesSecond        int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                    int    `json:"qosLcs,omitempty"`
	QosPacketRate             int    `json:"qosPacketRate,omitempty"`
	QosRtt                    int    `json:"qosRtt,omitempty"`
	QosTopology               int    `json:"qosTopology,omitempty"`
	QosVsCapacity             int    `json:"qosVsCapacity,omitempty"`
	QosVsScore                int    `json:"qosVsScore,omitempty"`
	TTL                       int    `json:"ttl,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
	MembersReference          struct {
		Link            string `json:"link,omitempty"`
		IsSubcollection bool   `json:"isSubcollection,omitempty"`
	}
}

// GetGTMAPools returns a list of all Pool/A records
func (b *BigIP) GetGTMAPools() (*GTMAPools, error) {
	var p GTMAPools
	err, _ := b.getForEntity(&p, uriGtm, uriPool, string(ARecord))
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetGTMAPool get's a Pool/A by name
func (b *BigIP) GetGTMAPool(name string) (*GTMAPool, error) {
	var w GTMAPool

	err, ok := b.getForEntity(&w, uriGtm, uriPool, string(ARecord), name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}

// AddGTMAPool adds a Pool/A by config to the BIG-IP system.
func (b *BigIP) AddGTMAPool(config *GTMAPool) error {
	return b.post(config, uriGtm, uriPool, string(ARecord))
}

// ModifyGTMAPool adds a Pool/A by config to the BIG-IP system.
func (b *BigIP) ModifyGTMAPool(fullPath string, config *GTMAPool) error {
	return b.put(config, uriGtm, uriPool, string(ARecord), fullPath)
}

// ********************************************************************************************************************
// *****************************************                        ***************************************************
// *****************************************   GTM A Pool Members   ***************************************************
// *****************************************                        ***************************************************
// ********************************************************************************************************************

// GTMAPoolMembers contains a list of every gtm/pool/a/members on the BIG-IP system.
type GTMAPoolMembers struct {
	GTMAPoolMembers []GTMAPoolMember `json:"items"`
}

// GTMAPoolMember contains information about each gtm/pool/a
type GTMAPoolMember struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	SubPath                   string `json:"subPath,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AppService                string `json:"appService,omitempty"`
	Description               string `json:"description,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	LimitMaxBps               uint64 `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       uint64 `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               uint64 `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	MemberOrder               int    `json:"memberOrder,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	Ratio                     int    `json:"ratio,omitempty"`
}

func buildPoolMemberFullPath(serverFullPath, poolMemberFullPath string) string {
	return fmt.Sprintf("%s:%s", serverFullPath, poolMemberFullPath)
}

// GetGTMAPoolMembers returns a list of all Pool/A Members records
func (b *BigIP) GetGTMAPoolMembers(fullPathToAPool string) (*GTMAPoolMembers, error) {
	var m GTMAPoolMembers
	err, _ := b.getForEntity(&m, uriGtm, uriPool, string(ARecord), fullPathToAPool, uriPoolMembers)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetGTMAPoolMember get's a Pool/A Member by name
func (b *BigIP) GetGTMAPoolMember(fullPathToAPool, serverFullPath, poolMemberFullPath string) (*GTMAPoolMember, error) {
	var w GTMAPoolMember

	fullPathToPoolMember := buildPoolMemberFullPath(serverFullPath, poolMemberFullPath)

	err, ok := b.getForEntity(&w, uriGtm, uriPool, string(ARecord), fullPathToAPool, uriPoolMember, fullPathToPoolMember)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}

// CreateGTMAPoolMember adds a Pool/A Member by using Paths, helpfull if Virtual Server Discovery is turned on
func (b *BigIP) CreateGTMAPoolMember(fullPathToAPool, serverFullPath, poolMemberFullPath string) error {
	config := &GTMAPoolMember{}
	config.Name = buildPoolMemberFullPath(serverFullPath, poolMemberFullPath)
	return b.post(config, uriGtm, uriPool, string(ARecord), fullPathToAPool, uriPoolMember)
}

// DeleteGTMAPoolMember remvoes a Pool/A Member
func (b *BigIP) DeleteGTMAPoolMember(fullPathToAPool, serverFullPath, poolMemberFullPath string) error {
	fullPathToPoolMember := buildPoolMemberFullPath(serverFullPath, poolMemberFullPath)
	return b.delete(uriGtm, uriPool, string(ARecord), fullPathToAPool, uriPoolMember, fullPathToPoolMember)
}

// GTMCNamePools contains a list of every gtm/pool/cname on the BIG-IP system.
type GTMCNamePools struct {
	GTMCNamePools []GTMCNamePool `json:"items"`
}

// GTMCNamePool contains information about each gtm/pool/cname.
type GTMCNamePool struct {
	Name                     string `json:"name,omitempty"`
	Partition                string `json:"partition,omitempty"`
	FullPath                 string `json:"fullPath,omitempty"`
	Generation               int    `json:"generation,omitempty"`
	AppService               string `json:"appService,omitempty"`
	Description              string `json:"description,omitempty"`
	Disabled                 bool   `json:"disabled,omitempty"`
	DynamicRatio             string `json:"dynamicRatio,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
	FallbackMode             string `json:"fallbackMode,omitempty"`
	LoadBalancingMode        string `json:"loadBalancingMode,omitempty"`
	ManualResume             string `json:"manualResume,omitempty"`
	TmPartition              string `json:"tmPartition,omitempty"`
	QosHitRatio              int    `json:"qosHitRatio,omitempty"`
	QosHops                  int    `json:"qosHops,omitempty"`
	QosKilobytesSecond       int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                   int    `json:"qosLcs,omitempty"`
	QosPacketRate            int    `json:"qosPacketRate,omitempty"`
	QosRtt                   int    `json:"qosRtt,omitempty"`
	QosTopology              int    `json:"qosTopology,omitempty"`
	QosVsCapacity            int    `json:"qosVsCapacity,omitempty"`
	QosVsScore               int    `json:"qosVsScore,omitempty"`
	TTL                      int    `json:"ttl,omitempty"`
	VerifyMemberAvailability string `json:"verifyMemberAvailability,omitempty"`
	MembersReference         struct {
		Link            string `json:"link,omitempty"`
		IsSubcollection bool   `json:"isSubcollection,omitempty"`
	}
}

// GetGTMCNamePools returns a list of all Pool/CNAME records.
func (b *BigIP) GetGTMCNamePools() (*GTMCNamePools, error) {
	var p GTMCNamePools
	err, _ := b.getForEntity(&p, uriGtm, uriPool, string(CNAMERecord))
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetGTMCNamePool gets a Pool/CNAME by name.
func (b *BigIP) GetGTMCNamePool(name string) (*GTMCNamePool, error) {
	var w GTMCNamePool

	err, ok := b.getForEntity(&w, uriGtm, uriPool, string(CNAMERecord), name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}

// GTMCNamePoolMembers contains a list of every gtm/pool/cname/members on the BIG-IP system.
type GTMCNamePoolMembers struct {
	GTMCNamePoolMembers []GTMCNamePoolMember `json:"items"`
}

// GTMCNamePoolMember contains information about each gtm/pool/cname member.
type GTMCNamePoolMember struct {
	Name         string `json:"name,omitempty"`
	FullPath     string `json:"fullPath,omitempty"`
	Generation   int    `json:"generation,omitempty"`
	SelfLink     string `json:"selfLink,omitempty"`
	Enabled      bool   `json:"enabled,omitempty"`
	MemberOrder  int    `json:"memberOrder,omitempty"`
	Ratio        int    `json:"ratio,omitempty"`
	StaticTarget string `json:"staticTarget,omitempty"`
}

// GetGTMCNamePoolMembers returns a list of all Pool/CName member records.
func (b *BigIP) GetGTMCNamePoolMembers(fullPathToCNamePool string) (*GTMCNamePoolMembers, error) {
	var m GTMCNamePoolMembers
	err, _ := b.getForEntity(&m, uriGtm, uriPool, string(CNAMERecord), fullPathToCNamePool, uriPoolMembers)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetGTMCNamePoolMember gets a Pool/CNAME member by name.
func (b *BigIP) GetGTMCNamePoolMember(fullPathToAPool, poolMemberFullPath string) (*GTMCNamePoolMember, error) {
	var w GTMCNamePoolMember

	err, ok := b.getForEntity(&w, uriGtm, uriPool, string(CNAMERecord), fullPathToAPool, uriPoolMember, poolMemberFullPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &w, nil
}

/*

These are here for later use -- so no one has to do this painful work!!!



// GTMAAAAPools contains a list of every gtm/pool/aaaa on the BIG-IP system.
type GTMAAAAPools struct {
	GTMAAAAPools []GTMAAAAPool `json:"items"`
}

// GTMAAAAPool contains information about each gtm/pool/aaaa
type GTMAAAAPool struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
	FullPath                  string `json:"fullPath,omitempty"`
	Generation                int    `json:"generation,omitempty"`
	AppService                string `json:"appService,omitempty"`
	Description               string `json:"description,omitempty"`
	Disabled                  bool   `json:"disabled,omitempty"`
	DynamicRatio              string `json:"dynamicRatio,omitempty"`
	Enabled                   bool   `json:"enabled,omitempty"`
	FallbackIP                string `json:"fallbackIp,omitempty"`
	FallbackMode              string `json:"fallbackMode,omitempty"`
	LimitMaxBps               uint64 `json:"limitMaxBps,omitempty"`
	LimitMaxBpsStatus         string `json:"limitMaxBpsStatus,omitempty"`
	LimitMaxConnections       uint64 `json:"limitMaxConnections,omitempty"`
	LimitMaxConnectionsStatus string `json:"limitMaxConnectionsStatus,omitempty"`
	LimitMaxPps               uint64 `json:"limitMaxPps,omitempty"`
	LimitMaxPpsStatus         string `json:"limitMaxPpsStatus,omitempty"`
	LoadBalancingMode         string `json:"loadBalancingMode,omitempty"`
	ManualResume              string `json:"manualResume,omitempty"`
	MaxAnswersReturned        int    `json:"maxAnswersReturned,omitempty"`
	Monitor                   string `json:"monitor,omitempty"`
	TmPartition               string `json:"tmPartition,omitempty"`
	QosHitRatio               int    `json:"qosHitRatio,omitempty"`
	QosHops                   int    `json:"qosHops,omitempty"`
	QosKilobytesSecond        int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                    int    `json:"qosLcs,omitempty"`
	QosPacketRate             int    `json:"qosPacketRate,omitempty"`
	QosRtt                    int    `json:"qosRtt,omitempty"`
	QosTopology               int    `json:"qosTopology,omitempty"`
	QosVsCapacity             int    `json:"qosVsCapacity,omitempty"`
	QosVsScore                int    `json:"qosVsScore,omitempty"`
	TTL                       int    `json:"ttl,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
}

// GTMMXPools contains a list of every gtm/pool/mx on the BIG-IP system.
type GTMMXPools struct {
	GTMMXPools []GTMMXPool `json:"items"`
}

// GTMMXPool contains information about each gtm/pool/mx
type GTMMXPool struct {
	Name                     string `json:"name,omitempty"`
	Partition                string `json:"partition,omitempty"`
	FullPath                 string `json:"fullPath,omitempty"`
	Generation               int    `json:"generation,omitempty"`
	AppService               string `json:"appService,omitempty"`
	Description              string `json:"description,omitempty"`
	Disabled                 bool   `json:"disabled,omitempty"`
	DynamicRatio             string `json:"dynamicRatio,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
	FallbackMode             string `json:"fallbackMode,omitempty"`
	LoadBalancingMode        string `json:"loadBalancingMode,omitempty"`
	ManualResume             string `json:"manualResume,omitempty"`
	TmPartition              string `json:"tmPartition,omitempty"`
	QosHitRatio              int    `json:"qosHitRatio,omitempty"`
	QosHops                  int    `json:"qosHops,omitempty"`
	QosKilobytesSecond       int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                   int    `json:"qosLcs,omitempty"`
	QosPacketRate            int    `json:"qosPacketRate,omitempty"`
	QosRtt                   int    `json:"qosRtt,omitempty"`
	QosTopology              int    `json:"qosTopology,omitempty"`
	QosVsCapacity            int    `json:"qosVsCapacity,omitempty"`
	QosVsScore               int    `json:"qosVsScore,omitempty"`
	TTL                      int    `json:"ttl,omitempty"`
	VerifyMemberAvailability string `json:"verifyMemberAvailability,omitempty"`
}

// GTMSrvPools contains a list of every gtm/pool/srv on the BIG-IP system.
type GTMSrvPools struct {
	GTMSrvPools []GTMSrvPool `json:"items"`
}

// GTMSrvPool contains information about each gtm/pool/srv
type GTMSrvPool struct {
	Name                     string `json:"name,omitempty"`
	Partition                string `json:"partition,omitempty"`
	FullPath                 string `json:"fullPath,omitempty"`
	Generation               int    `json:"generation,omitempty"`
	AppService               string `json:"appService,omitempty"`
	Description              string `json:"description,omitempty"`
	Disabled                 bool   `json:"disabled,omitempty"`
	DynamicRatio             string `json:"dynamicRatio,omitempty"`
	Enabled                  bool   `json:"enabled,omitempty"`
	FallbackMode             string `json:"fallbackMode,omitempty"`
	LoadBalancingMode        string `json:"loadBalancingMode,omitempty"`
	ManualResume             string `json:"manualResume,omitempty"`
	TmPartition              string `json:"tmPartition,omitempty"`
	QosHitRatio              int    `json:"qosHitRatio,omitempty"`
	QosHops                  int    `json:"qosHops,omitempty"`
	QosKilobytesSecond       int    `json:"qosKilobytesSecond,omitempty"`
	QosLcs                   int    `json:"qosLcs,omitempty"`
	QosPacketRate            int    `json:"qosPacketRate,omitempty"`
	QosRtt                   int    `json:"qosRtt,omitempty"`
	QosTopology              int    `json:"qosTopology,omitempty"`
	QosVsCapacity            int    `json:"qosVsCapacity,omitempty"`
	QosVsScore               int    `json:"qosVsScore,omitempty"`
	TTL                      int    `json:"ttl,omitempty"`
	VerifyMemberAvailability string `json:"verifyMemberAvailability,omitempty"`
}

*/
