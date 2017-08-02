package bigip

// FastL4Profile is a representation of a fastL4Profile configuration
type FastL4Profile struct {
	Kind                    string `json:"kind"`
	Name                    string `json:"name"`
	Partition               string `json:"partition"`
	FullPath                string `json:"fullPath"`
	Generation              int    `json:"generation"`
	SelfLink                string `json:"selfLink"`
	HardwareSynCookie       string `json:"hardwareSynCookie"`
	IdleTimeout             string `json:"idleTimeout"`
	IPTosToClient           string `json:"ipTosToClient"`
	IPTosToServer           string `json:"ipTosToServer"`
	KeepAliveInterval       string `json:"keepAliveInterval"`
	LinkQosToClient         string `json:"linkQosToClient"`
	LinkQosToServer         string `json:"linkQosToServer"`
	LooseClose              string `json:"looseClose"`
	LooseInitialization     string `json:"looseInitialization"`
	MssOverride             int    `json:"mssOverride"`
	PriorityToClient        string `json:"priorityToClient"`
	PriorityToServer        string `json:"priorityToServer"`
	PvaAcceleration         string `json:"pvaAcceleration"`
	PvaDynamicClientPackets int    `json:"pvaDynamicClientPackets"`
	PvaDynamicServerPackets int    `json:"pvaDynamicServerPackets"`
	PvaFlowAging            string `json:"pvaFlowAging"`
	PvaFlowEvict            string `json:"pvaFlowEvict"`
	PvaOffloadDynamic       string `json:"pvaOffloadDynamic"`
	PvaOffloadState         string `json:"pvaOffloadState"`
	ReassembleFragments     string `json:"reassembleFragments"`
	ReceiveWindowSize       int    `json:"receiveWindowSize"`
	ResetOnTimeout          string `json:"resetOnTimeout"`
	RttFromClient           string `json:"rttFromClient"`
	RttFromServer           string `json:"rttFromServer"`
	ServerSack              string `json:"serverSack"`
	ServerTimestamp         string `json:"serverTimestamp"`
	SoftwareSynCookie       string `json:"softwareSynCookie"`
	TCPCloseTimeout         string `json:"tcpCloseTimeout"`
	TCPGenerateIsn          string `json:"tcpGenerateIsn"`
	TCPHandshakeTimeout     string `json:"tcpHandshakeTimeout"`
	TCPStripSack            string `json:"tcpStripSack"`
	TCPTimestampMode        string `json:"tcpTimestampMode"`
	TCPWscaleMode           string `json:"tcpWscaleMode"`
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
