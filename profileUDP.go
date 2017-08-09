package bigip

// UDPProfile represents a UDP Profile configuration
type UDPProfile struct {
	Kind                  string `json:"kind,omitempty"`
	Name                  string `json:"name,omitempty"`
	Partition             string `json:"partition,omitempty"`
	FullPath              string `json:"fullPath,omitempty"`
	Generation            int    `json:"generation,omitempty"`
	SelfLink              string `json:"selfLink,omitempty"`
	AllowNoPayload        string `json:"allowNoPayload,omitempty"`
	DatagramLoadBalancing string `json:"datagramLoadBalancing,omitempty"`
	DefaultsFrom          string `json:"defaultsFrom,omitempty"`
	IdleTimeout           string `json:"idleTimeout,omitempty"`
	IPTosToClient         string `json:"ipTosToClient,omitempty"`
	LinkQosToClient       string `json:"linkQosToClient,omitempty"`
	NoChecksum            string `json:"noChecksum,omitempty"`
	ProxyMss              string `json:"proxyMss,omitempty"`
}

// UDPProfiles is an array of UDPProfile structs
type UDPProfiles []UDPProfile

// UDPProfiles returns a list of udp profiles.
func (b *BigIP) UDPProfiles() (*UDPProfiles, error) {
	var serverUDPProfiles UDPProfiles
	err, _ := b.getForEntity(&serverUDPProfiles, uriLtm, uriProfile, uriProfileUDP)
	if err != nil {
		return nil, err
	}

	return &serverUDPProfiles, nil
}

// GetUDPProfile gets a udp profile by name. Returns nil if the udp profile does not exist
func (b *BigIP) GetUDPProfile(name string) (*UDPProfile, error) {
	var serverUDPProfile UDPProfile
	err, ok := b.getForEntity(&serverUDPProfile, uriLtm, uriProfile, uriProfileUDP, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverUDPProfile, nil
}

// CreateUDPProfile creates a new udp profile on the BIG-IP system.
func (b *BigIP) CreateUDPProfile(name string, parent string) error {
	config := &UDPProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriProfileUDP)
}

// AddUDPProfile adds a new udp profile on the BIG-IP system.
func (b *BigIP) AddUDPProfile(config *UDPProfile) error {
	return b.post(config, uriLtm, uriProfile, uriProfileUDP)
}

// DeleteUDPProfile removes a udp profile.
func (b *BigIP) DeleteUDPProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriProfileUDP, name)
}

// ModifyUDPProfile allows you to change any attribute of a udp profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyUDPProfile(name string, config *UDPProfile) error {
	return b.put(config, uriLtm, uriProfile, uriProfileUDP, name)
}
