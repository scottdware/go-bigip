package bigip

// GTM Documentation
// https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_gtm.ashx

// **********************************
// **   GTM WideIp
// **********************************

// GTMWideIPs contains a list of every WideIP on the BIG-IP system.
type GTMWideIPs struct {
	GTMWideIPs []GTMWideIP `json:"items"`
}

// GTMWideIP contains information about each wide ip (regardless of type: A, AAAA, CNAME, etc.)
// Type is what determine the type of record the WideIp is for in the docs, however that is NOT returned by the API
// Instead you have to query the Type by the uri   wideip/a  wideip/cname  that = type
type GTMWideIP struct {
	Name                              string `json:"name,omitempty"`
	Partition                         string `json:"partition,omitempty"`
	FullPath                          string `json:"fullPath,omitempty"`
	Generation                        int    `json:"generation,omitempty"`
	AppService                        string `json:"appService,omitempty"`
	Description                       string `json:"description,omitempty"`
	Disabled                          bool   `json:"disabled,omitempty"`
	Enabled                           bool   `json:"enabled,omitempty"`
	FailureRcode                      string `json:"failureRcode,omitempty"`
	FailureRcodeResponse              string `json:"failureRcodeResponse,omitempty"`
	FailureRcodeTTL                   int    `json:"failureRcodeTtl,omitempty"`
	LastResortPool                    string `json:"lastResortPool,omitempty"`
	LoadBalancingDecisionLogVerbosity string `json:"loadBalancingDecisionLogVerbosity,omitempty"`
	MinimalResponse                   string `json:"minimalResponse,omitempty"`
	PersistCidrIpv4                   int    `json:"persistCidrIpv4,omitempty"`
	PersistCidrIpv6                   int    `json:"persistCidrIpv6,omitempty"`
	Persistence                       string `json:"persistence,omitempty"`
	PoolLbMode                        string `json:"poolLbMode,omitempty"`
	TTLPersistence                    int    `json:"ttlPersistence,omitempty"`

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

// GTMWideIPs returns a list of all WideIps for a provided type
func (b *BigIP) GTMWideIPs(recordType GTMType) (*GTMWideIPs, error) {
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

// **********************************
// **   GTM Pool A
// **********************************

/*

These are here for later use -- so no one has to do this painful work!!!


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
	Type                      string `json:"type,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
	MembersReference          struct {
		Link            string `json:"link,omitempty"`
		IsSubcollection bool   `json:"isSubcollection,omitempty"`
	}
}

// GTMAPoolMembers contains a list of every gtm/pool/a/members on the BIG-IP system.
type GTMAPoolMembers struct {
	GTMAPoolMembers []GTMAPoolMember `json:"itmes"`
}

// GTMAPoolMember contains information about each gtm/pool/a
type GTMAPoolMember struct {
	Name                      string `json:"name,omitempty"`
	Partition                 string `json:"partition,omitempty"`
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
	Type                      string `json:"type,omitempty"`
}

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
	Type                      string `json:"type,omitempty"`
	VerifyMemberAvailability  string `json:"verifyMemberAvailability,omitempty"`
}

// GTMCNamePools contains a list of every gtm/pool/cname on the BIG-IP system.
type GTMCNamePools struct {
	GTMCNamePools []GTMCNamePool `json:"items"`
}

// GTMCNamePool contains information about each gtm/pool/cname
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
	Type                     string `json:"type,omitempty"`
	VerifyMemberAvailability string `json:"verifyMemberAvailability,omitempty"`
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
	Type                     string `json:"type,omitempty"`
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
	Type                     string `json:"type,omitempty"`
	VerifyMemberAvailability string `json:"verifyMemberAvailability,omitempty"`
}

*/
