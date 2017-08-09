package bigip

// OneConnectProfiles
// Documentation: https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_ltm_profile_oneConnect.ashx

// OneConnectProfiles contains a list of every oneConnect profile on the BIG-IP system.
type OneConnectProfiles struct {
	OneConnectProfiles []OneConnectProfile `json:"items"`
}

// OneConnectProfile contains information about each oneConnect profile. You can use all
// of these fields when modifying a oneConnect profile.
type OneConnectProfile struct {
	Kind                string `json:"kind,omitempty"`
	Name                string `json:"name,omitempty"`
	Partition           string `json:"partition,omitempty"`
	FullPath            string `json:"fullPath,omitempty"`
	Generation          int    `json:"generation,omitempty"`
	SelfLink            string `json:"selfLink,omitempty"`
	IdleTimeoutOverride string `json:"idleTimeoutOverride,omitempty"`
	MaxAge              int    `json:"maxAge,omitempty"`
	MaxReuse            int    `json:"maxReuse,omitempty"`
	MaxSize             int    `json:"maxSize,omitempty"`
	SharePools          string `json:"sharePools,omitempty"`
	SourceMask          string `json:"sourceMask,omitempty"`
	DefaultsFrom        string `json:"defaultsFrom,omitempty"`
}

// OneConnectProfiles returns a list of oneConnect profiles.
func (b *BigIP) OneConnectProfiles() (*OneConnectProfiles, error) {
	var oneConnectProfiles OneConnectProfiles
	err, _ := b.getForEntity(&oneConnectProfiles, uriLtm, uriProfile, uriOneConnect)
	if err != nil {
		return nil, err
	}

	return &oneConnectProfiles, nil
}

// GetOneConnectProfile gets a oneConnect profile by name. Returns nil if the oneConnect profile does not exist
func (b *BigIP) GetOneConnectProfile(name string) (*OneConnectProfile, error) {
	var oneConnectProfile OneConnectProfile
	err, ok := b.getForEntity(&oneConnectProfile, uriLtm, uriProfile, uriOneConnect, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &oneConnectProfile, nil
}

// CreateOneConnectProfile creates a new oneConnect profile on the BIG-IP system.
func (b *BigIP) CreateOneConnectProfile(name string, parent string) error {
	config := &OneConnectProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriOneConnect)
}

// AddOneConnectProfile adds a new oneConnect profile on the BIG-IP system.
func (b *BigIP) AddOneConnectProfile(config *OneConnectProfile) error {
	return b.post(config, uriLtm, uriProfile, uriOneConnect)
}

// DeleteOneConnectProfile removes a oneConnect profile.
func (b *BigIP) DeleteOneConnectProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriOneConnect, name)
}

// ModifyOneConnectProfile allows you to change any attribute of a sever-oneConnect profile.
// Fields that can be modified are referenced in the VirtualClient struct.
func (b *BigIP) ModifyOneConnectProfile(name string, config *OneConnectProfile) error {
	return b.put(config, uriLtm, uriProfile, uriOneConnect, name)
}
