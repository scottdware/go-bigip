package bigip

// FTPProfile represents a FTP Profile configuration
type FTPProfile struct {
	Kind                 string `json:"kind"`
	Name                 string `json:"name"`
	Partition            string `json:"partition,omitempty"`
	FullPath             string `json:"fullPath,omitempty"`
	Generation           int    `json:"generation,omitempty"`
	SelfLink             string `json:"selfLink,omitempty"`
	InheritParentProfile string `json:"inheritParentProfile,omitempty"`
	Port                 int    `json:"port,omitempty"`
	Security             string `json:"security,omitempty"`
	TranslateExtended    string `json:"translateExtended,omitempty"`
	DefaultsFrom         string `json:defaultsFrom`
}

// FTPProfiles is an array of FTPProfile structs
type FTPProfiles []FTPProfile

// FTPProfiles returns a list of FTP profiles.
func (b *BigIP) FTPProfiles() (*FTPProfiles, error) {
	var serverFTPProfiles FTPProfiles
	err, _ := b.getForEntity(&serverFTPProfiles, uriLtm, uriProfile, uriProfileFTP)
	if err != nil {
		return nil, err
	}

	return &serverFTPProfiles, nil
}

// GetFTPProfile gets a FTP profile by name. Returns nil if the FTP profile does not exist
func (b *BigIP) GetFTPProfile(name string) (*FTPProfile, error) {
	var serverFTPProfile FTPProfile
	err, ok := b.getForEntity(&serverFTPProfile, uriLtm, uriProfile, uriProfileFTP, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverFTPProfile, nil
}

// CreateFTPProfile creates a new FTP profile on the BIG-IP system.
func (b *BigIP) CreateFTPProfile(name string) error {
	config := &FTPProfile{
		Name: name,
	}

	return b.post(config, uriLtm, uriProfile, uriProfileFTP)
}

// AddFTPProfile adds a new FTP profile on the BIG-IP system.
func (b *BigIP) AddFTPProfile(config *FTPProfile) error {
	return b.post(config, uriLtm, uriProfile, uriProfileFTP)
}

// DeleteFTPProfile removes a FTP profile.
func (b *BigIP) DeleteFTPProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriProfileFTP, name)
}

// ModifyFTPProfile allows you to change any attribute of a FTP profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyFTPProfile(name string, config *FTPProfile) error {
	return b.put(config, uriLtm, uriProfile, uriProfileFTP, name)
}
