package bigip

// TCPProfile represents a TCP Profile configuration
type TCPProfile struct {
	Kind                string `json:"kind,omitempty"`
	Name                string `json:"name,omitempty"`
	Partition           string `json:"partition,omitempty"`
	FullPath            string `json:"fullPath,omitempty"`
	Generation          int    `json:"generation,omitempty"`
	SelfLink            string `json:"selfLink,omitempty"`
	Abc                 string `json:"abc,omitempty"`
	AckOnPush           string `json:"ackOnPush,omitempty"`
	CloseWaitTimeout    int    `json:"closeWaitTimeout,omitempty"`
	CmetricsCache       string `json:"cmetricsCache,omitempty"`
	CongestionControl   string `json:"congestionControl,omitempty"`
	DefaultsFrom        string `json:"defaultsFrom,omitempty"`
	DeferredAccept      string `json:"deferredAccept,omitempty"`
	DelayWindowControl  string `json:"delayWindowControl,omitempty"`
	DelayedAcks         string `json:"delayedAcks,omitempty"`
	Dsack               string `json:"dsack,omitempty"`
	Ecn                 string `json:"ecn,omitempty"`
	FinWaitTimeout      int    `json:"finWaitTimeout,omitempty"`
	HardwareSynCookie   string `json:"hardwareSynCookie,omitempty"`
	IdleTimeout         int    `json:"idleTimeout,omitempty"`
	InitCwnd            int    `json:"initCwnd,omitempty"`
	InitRwnd            int    `json:"initRwnd,omitempty"`
	IPTosToClient       string `json:"ipTosToClient,omitempty"`
	KeepAliveInterval   int    `json:"keepAliveInterval,omitempty"`
	LimitedTransmit     string `json:"limitedTransmit,omitempty"`
	LinkQosToClient     string `json:"linkQosToClient,omitempty"`
	MaxRetrans          int    `json:"maxRetrans,omitempty"`
	MaxSegmentSize      int    `json:"maxSegmentSize,omitempty"`
	Md5Signature        string `json:"md5Signature,omitempty"`
	MinimumRto          int    `json:"minimumRto,omitempty"`
	Mptcp               string `json:"mptcp,omitempty"`
	MptcpCsum           string `json:"mptcpCsum,omitempty"`
	MptcpCsumVerify     string `json:"mptcpCsumVerify,omitempty"`
	MptcpDebug          string `json:"mptcpDebug,omitempty"`
	MptcpFallback       string `json:"mptcpFallback,omitempty"`
	MptcpFastjoin       string `json:"mptcpFastjoin,omitempty"`
	MptcpJoinMax        int    `json:"mptcpJoinMax,omitempty"`
	MptcpMakeafterbreak string `json:"mptcpMakeafterbreak,omitempty"`
	MptcpNojoindssack   string `json:"mptcpNojoindssack,omitempty"`
	MptcpRtomax         int    `json:"mptcpRtomax,omitempty"`
	MptcpRxmitmin       int    `json:"mptcpRxmitmin,omitempty"`
	MptcpSubflowmax     int    `json:"mptcpSubflowmax,omitempty"`
	MptcpTimeout        int    `json:"mptcpTimeout,omitempty"`
	Nagle               string `json:"nagle,omitempty"`
	PktLossIgnoreBurst  int    `json:"pktLossIgnoreBurst,omitempty"`
	PktLossIgnoreRate   int    `json:"pktLossIgnoreRate,omitempty"`
	ProxyBufferHigh     int    `json:"proxyBufferHigh,omitempty"`
	ProxyBufferLow      int    `json:"proxyBufferLow,omitempty"`
	ProxyMss            string `json:"proxyMss,omitempty"`
	ProxyOptions        string `json:"proxyOptions,omitempty"`
	RatePace            string `json:"ratePace,omitempty"`
	ReceiveWindowSize   int    `json:"receiveWindowSize,omitempty"`
	ResetOnTimeout      string `json:"resetOnTimeout,omitempty"`
	SelectiveAcks       string `json:"selectiveAcks,omitempty"`
	SelectiveNack       string `json:"selectiveNack,omitempty"`
	SendBufferSize      int    `json:"sendBufferSize,omitempty"`
	SlowStart           string `json:"slowStart,omitempty"`
	SynMaxRetrans       int    `json:"synMaxRetrans,omitempty"`
	SynRtoBase          int    `json:"synRtoBase,omitempty"`
	TimeWaitRecycle     string `json:"timeWaitRecycle,omitempty"`
	TimeWaitTimeout     int    `json:"timeWaitTimeout,omitempty"`
	Timestamps          string `json:"timestamps,omitempty"`
	VerifiedAccept      string `json:"verifiedAccept,omitempty"`
	ZeroWindowTimeout   int    `json:"zeroWindowTimeout,omitempty"`
}

// TCPProfiles is an array of TCPProfile structs
type TCPProfiles []TCPProfile

// TCPProfiles returns a list of tcp profiles.
func (b *BigIP) TCPProfiles() (*TCPProfiles, error) {
	var serverTCPProfiles TCPProfiles
	err, _ := b.getForEntity(&serverTCPProfiles, uriLtm, uriProfile, uriProfileTCP)
	if err != nil {
		return nil, err
	}

	return &serverTCPProfiles, nil
}

// GetTCPProfile gets a tcp profile by name. Returns nil if the tcp profile does not exist
func (b *BigIP) GetTCPProfile(name string) (*TCPProfile, error) {
	var serverTCPProfile TCPProfile
	err, ok := b.getForEntity(&serverTCPProfile, uriLtm, uriProfile, uriProfileTCP, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverTCPProfile, nil
}

// CreateTCPProfile creates a new tcp profile on the BIG-IP system.
func (b *BigIP) CreateTCPProfile(name string, parent string) error {
	config := &TCPProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriProfileTCP)
}

// AddTCPProfile adds a new tcp profile on the BIG-IP system.
func (b *BigIP) AddTCPProfile(config *TCPProfile) error {
	return b.post(config, uriLtm, uriProfile, uriProfileTCP)
}

// DeleteTCPProfile removes a tcp profile.
func (b *BigIP) DeleteTCPProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriProfileTCP, name)
}

// ModifyTCPProfile allows you to change any attribute of a tcp profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyTCPProfile(name string, config *TCPProfile) error {
	return b.put(config, uriLtm, uriProfile, uriProfileTCP, name)
}
