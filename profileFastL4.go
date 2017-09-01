package bigip

// FastL4Profile is a representation of a fastL4Profile configuration
type FastL4Profile struct {
	Kind                    string `json:"kind,omitempty"`
	Name                    string `json:"name"`
	Partition               string `json:"partition,omitempty"`
	FullPath                string `json:"fullPath,omitempty"`
	Generation              int    `json:"generation,omitempty"`
	SelfLink                string `json:"selfLink,omitempty"`
	HardwareSynCookie       string `json:"hardwareSynCookie,omitempty"`
	IdleTimeout             string `json:"idleTimeout,omitempty"`
	IPTosToClient           string `json:"ipTosToClient,omitempty"`
	IPTosToServer           string `json:"ipTosToServer,omitempty"`
	KeepAliveInterval       string `json:"keepAliveInterval,omitempty"`
	LinkQosToClient         string `json:"linkQosToClient,omitempty"`
	LinkQosToServer         string `json:"linkQosToServer,omitempty"`
	LooseClose              string `json:"looseClose,omitempty"`
	LooseInitialization     string `json:"looseInitialization,omitempty"`
	MssOverride             int    `json:"mssOverride,omitempty"`
	PriorityToClient        int    `json:"priorityToClient,omitempty"`
	PriorityToServer        int    `json:"priorityToServer,omitempty"`
	PvaAcceleration         string `json:"pvaAcceleration,omitempty"`
	PvaDynamicClientPackets int    `json:"pvaDynamicClientPackets,omitempty"`
	PvaDynamicServerPackets int    `json:"pvaDynamicServerPackets,omitempty"`
	PvaFlowAging            string `json:"pvaFlowAging,omitempty"`
	PvaFlowEvict            string `json:"pvaFlowEvict,omitempty"`
	PvaOffloadDynamic       string `json:"pvaOffloadDynamic,omitempty"`
	PvaOffloadState         string `json:"pvaOffloadState,omitempty"`
	ReassembleFragments     string `json:"reassembleFragments,omitempty"`
	ReceiveWindowSize       int    `json:"receiveWindowSize,omitempty"`
	ResetOnTimeout          string `json:"resetOnTimeout,omitempty"`
	RttFromClient           string `json:"rttFromClient,omitempty"`
	RttFromServer           string `json:"rttFromServer,omitempty"`
	ServerSack              string `json:"serverSack,omitempty"`
	ServerTimestamp         string `json:"serverTimestamp,omitempty"`
	SoftwareSynCookie       string `json:"softwareSynCookie,omitempty"`
	TCPCloseTimeout         string `json:"tcpCloseTimeout,omitempty"`
	TCPGenerateIsn          string `json:"tcpGenerateIsn,omitempty"`
	TCPHandshakeTimeout     string `json:"tcpHandshakeTimeout,omitempty"`
	TCPStripSack            string `json:"tcpStripSack,omitempty"`
	TCPTimestampMode        string `json:"tcpTimestampMode,omitempty"`
	TCPWscaleMode           string `json:"tcpWscaleMode,omitempty"`
	DefaultsFrom            string `json:"defaultsFrom"`
}

// FastL4Profiles is an array of FastL4Profile structs
type FastL4Profiles []FastL4Profile

// FastL4Profiles returns a list of fastL4 profiles.
func (b *BigIP) FastL4Profiles() (*FastL4Profiles, error) {
	var serverFastL4Profiles FastL4Profiles
	err, _ := b.getForEntity(&serverFastL4Profiles, uriLtm, uriProfile, uriProfileFastL4)
	if err != nil {
		return nil, err
	}

	return &serverFastL4Profiles, nil
}

// GetFastL4Profile gets a fastL4 profile by name. Returns nil if the fastL4 profile does not exist
func (b *BigIP) GetFastL4Profile(name string) (*FastL4Profile, error) {
	var serverFastL4Profile FastL4Profile
	err, ok := b.getForEntity(&serverFastL4Profile, uriLtm, uriProfile, uriProfileFastL4, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverFastL4Profile, nil
}

// CreateFastL4Profile creates a new fastL4 profile on the BIG-IP system.
func (b *BigIP) CreateFastL4Profile(name string, parent string) error {
	config := &FastL4Profile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriProfileFastL4)
}

// AddFastL4Profile adds a new fastL4 profile on the BIG-IP system.
func (b *BigIP) AddFastL4Profile(config *FastL4Profile) error {
	return b.post(config, uriLtm, uriProfile, uriProfileFastL4)
}

// DeleteFastL4Profile removes a fastL4 profile.
func (b *BigIP) DeleteFastL4Profile(name string) error {
	return b.delete(uriLtm, uriProfile, uriProfileFastL4, name)
}

// ModifyFastL4Profile allows you to change any attribute of a fastL4 profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyFastL4Profile(name string, config *FastL4Profile) error {
	return b.put(config, uriLtm, uriProfile, uriProfileFastL4, name)
}
