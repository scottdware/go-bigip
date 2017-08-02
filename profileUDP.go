package bigip

// UDPProfile represents a UDP Profile configuration
type UDPProfile struct {
	Kind                  string `json:"kind"`
	Name                  string `json:"name"`
	Partition             string `json:"partition"`
	FullPath              string `json:"fullPath"`
	Generation            int    `json:"generation"`
	SelfLink              string `json:"selfLink"`
	AllowNoPayload        string `json:"allowNoPayload"`
	DatagramLoadBalancing string `json:"datagramLoadBalancing"`
	DefaultsFrom          string `json:"defaultsFrom"`
	IdleTimeout           string `json:"idleTimeout"`
	IPTosToClient         string `json:"ipTosToClient"`
	LinkQosToClient       string `json:"linkQosToClient"`
	NoChecksum            string `json:"noChecksum"`
	ProxyMss              string `json:"proxyMss"`
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
