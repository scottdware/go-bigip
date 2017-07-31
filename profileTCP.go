package bigip

// TCPProfile represents a TCP Profile configuration
type TCPProfile struct {
	Kind                string `json:"kind"`
	Name                string `json:"name"`
	Partition           string `json:"partition"`
	FullPath            string `json:"fullPath"`
	Generation          int    `json:"generation"`
	SelfLink            string `json:"selfLink"`
	Abc                 string `json:"abc"`
	AckOnPush           string `json:"ackOnPush"`
	CloseWaitTimeout    int    `json:"closeWaitTimeout"`
	CmetricsCache       string `json:"cmetricsCache"`
	CongestionControl   string `json:"congestionControl"`
	DefaultsFrom        string `json:"defaultsFrom"`
	DeferredAccept      string `json:"deferredAccept"`
	DelayWindowControl  string `json:"delayWindowControl"`
	DelayedAcks         string `json:"delayedAcks"`
	Dsack               string `json:"dsack"`
	Ecn                 string `json:"ecn"`
	FinWaitTimeout      int    `json:"finWaitTimeout"`
	HardwareSynCookie   string `json:"hardwareSynCookie"`
	IdleTimeout         int    `json:"idleTimeout"`
	InitCwnd            int    `json:"initCwnd"`
	InitRwnd            int    `json:"initRwnd"`
	IPTosToClient       string `json:"ipTosToClient"`
	KeepAliveInterval   int    `json:"keepAliveInterval"`
	LimitedTransmit     string `json:"limitedTransmit"`
	LinkQosToClient     string `json:"linkQosToClient"`
	MaxRetrans          int    `json:"maxRetrans"`
	MaxSegmentSize      int    `json:"maxSegmentSize"`
	Md5Signature        string `json:"md5Signature"`
	MinimumRto          int    `json:"minimumRto"`
	Mptcp               string `json:"mptcp"`
	MptcpCsum           string `json:"mptcpCsum"`
	MptcpCsumVerify     string `json:"mptcpCsumVerify"`
	MptcpDebug          string `json:"mptcpDebug"`
	MptcpFallback       string `json:"mptcpFallback"`
	MptcpFastjoin       string `json:"mptcpFastjoin"`
	MptcpJoinMax        int    `json:"mptcpJoinMax"`
	MptcpMakeafterbreak string `json:"mptcpMakeafterbreak"`
	MptcpNojoindssack   string `json:"mptcpNojoindssack"`
	MptcpRtomax         int    `json:"mptcpRtomax"`
	MptcpRxmitmin       int    `json:"mptcpRxmitmin"`
	MptcpSubflowmax     int    `json:"mptcpSubflowmax"`
	MptcpTimeout        int    `json:"mptcpTimeout"`
	Nagle               string `json:"nagle"`
	PktLossIgnoreBurst  int    `json:"pktLossIgnoreBurst"`
	PktLossIgnoreRate   int    `json:"pktLossIgnoreRate"`
	ProxyBufferHigh     int    `json:"proxyBufferHigh"`
	ProxyBufferLow      int    `json:"proxyBufferLow"`
	ProxyMss            string `json:"proxyMss"`
	ProxyOptions        string `json:"proxyOptions"`
	RatePace            string `json:"ratePace"`
	ReceiveWindowSize   int    `json:"receiveWindowSize"`
	ResetOnTimeout      string `json:"resetOnTimeout"`
	SelectiveAcks       string `json:"selectiveAcks"`
	SelectiveNack       string `json:"selectiveNack"`
	SendBufferSize      int    `json:"sendBufferSize"`
	SlowStart           string `json:"slowStart"`
	SynMaxRetrans       int    `json:"synMaxRetrans"`
	SynRtoBase          int    `json:"synRtoBase"`
	TimeWaitRecycle     string `json:"timeWaitRecycle"`
	TimeWaitTimeout     int    `json:"timeWaitTimeout"`
	Timestamps          string `json:"timestamps"`
	VerifiedAccept      string `json:"verifiedAccept"`
	ZeroWindowTimeout   int    `json:"zeroWindowTimeout"`
}

// TCPProfiles is an array of TCPProfile structs
type TCPProfiles []TCPProfile

// TCPProfiles returns a list of tcp profiles.
func (b *BigIP) TCPProfiles() (*TCPProfiles, error) {
	var serverSSLProfiles TCPProfiles
	err, _ := b.getForEntity(&serverSSLProfiles, uriLtm, uriProfile, uriProfileTCP)
	if err != nil {
		return nil, err
	}

	return &serverSSLProfiles, nil
}

// GetTCPProfile gets a tcp profile by name. Returns nil if the tcp profile does not exist
func (b *BigIP) GetTCPProfile(name string) (*TCPProfile, error) {
	var serverSSLProfile TCPProfile
	err, ok := b.getForEntity(&serverSSLProfile, uriLtm, uriProfile, uriProfileTCP, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverSSLProfile, nil
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
