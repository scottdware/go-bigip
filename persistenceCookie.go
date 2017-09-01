package bigip

type PersistenceCookies struct {
	PersistenceCookie []PersistenceCookie `json:"items,omitempty"`
}

type PersistenceCookie struct {
	Name                       string `json:"name"`
	DefaultsFrom               string `json:"defaultsFrom"`
	Kind                       string `json:"kind,omitempty"`
	Mode                       string `json:"mode,omitempty"`
	Partition                  string `json:"partition,omitempty"`
	FullPath                   string `json:"fullPath,omitempty"`
	Generation                 int    `json:"generation,omitempty"`
	SelfLink                   string `json:"selfLink,omitempty"`
	AlwaysSend                 string `json:"alwaysSend,omitempty"`
	Description                string `json:"description,omitempty"`
	AppService                 string `json:"appService,omitempty"`
	CookieEncryption           string `json:"cookieEncryption,omitempty"`
	CookieEncryptionPassphrase string `json:"cookieEncryptionPassphrase,omitempty"`
	CookieName                 string `json:"cookieName,omitempty"`
	Expiration                 string `json:"expiration,omitempty"`
	HashLength                 int    `json:"hashLength,omitempty"`
	HashOffset                 int    `json:"hashOffset,omitempty"`
	MatchAcrossPools           string `json:"matchAcrossPools,omitempty"`
	MatchAcrossServices        string `json:"matchAcrossServices,omitempty"`
	MatchAcrossVirtuals        string `json:"matchAcrossVirtuals,omitempty"`
	Method                     string `json:"method,omitempty"`
	Mirror                     string `json:"mirror,omitempty"`
	Secure                     string `json:"secure,omitempty"`
	TMPartition                string `json:"tmPartition,omitempty"`
	OverrideConnectionLimit    string `json:"overrideConnectionLimit,omitempty"`
	Timeout                    string `json:"timeout,omitempty"`
}

// PersistenceCookie returns a list of oersistence profiles.
func (b *BigIP) PersistenceCookie() (*PersistenceCookie, error) {
	var persistenceProfiles PersistenceCookie
	err, _ := b.getForEntity(&persistenceProfiles, uriLtm, uriPersistences)
	if err != nil {
		return nil, err
	}

	return &persistenceProfiles, nil
}

// GetPersistenceCookie gets a persistence profile by name. Returns nil if the persistence profile does not exist
func (b *BigIP) GetPersistenceCookie(name string) (*PersistenceCookie, error) {
	var persistenceProfile PersistenceCookie
	err, ok := b.getForEntity(&persistenceProfile, uriLtm, uriProfile, uriPersistenceCookie, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &persistenceProfile, nil
}

// CreatePersistenceCookie creates a new persistence profile on the BIG-IP system.
func (b *BigIP) CreatePersistenceCookie(name string, parent string) error {
	config := &PersistenceCookie{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriPersistenceCookie)
}

// AddPersistenceCookie adds a new persistence profile on the BIG-IP system.
func (b *BigIP) AddPersistenceCookie(config *PersistenceCookie) error {
	return b.post(config, uriLtm, uriProfile, uriPersistenceCookie)
}

// DeletePersistenceCookie removes a persistence profile.
func (b *BigIP) DeletePersistenceCookie(name string) error {
	return b.delete(uriLtm, uriProfile, uriPersistenceCookie, name)
}

// ModifyPersistenceCookie allows you to change any attribute of a persistence profile.
// Fields that can be modified are referenced in the PersistenceCookie struct.
func (b *BigIP) ModifyPersistenceCookie(name string, config *PersistenceCookie) error {
	return b.put(config, uriLtm, uriProfile, uriPersistenceCookie, name)
}
