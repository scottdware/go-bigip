package bigip

type PersistenceCookies struct {
	PersistenceCookie []PersistenceCookie `json:"items"`
}

type PersistenceCookie struct {
	Name                       string `json:"name"`
	DefaultsFrom               string `json:"defaultsFrom"`
	Kind                       string `json:"kind"`
	Mode                       string `json:"mode"`
	Partition                  string `json:"partition"`
	FullPath                   string `json:"fullPath"`
	Generation                 int    `json:"generation"`
	SelfLink                   string `json:"selfLink"`
	AlwaysSend                 string `json:"alwaysSend"`
	Description                string `json:"description"`
	AppService                 string `json:"appService"`
	CookieEncryption           string `json:"cookieEncryption"`
	CookieEncryptionPassphrase string `json:"cookieEncryptionPassphrase"`
	CookieName                 string `json:"cookieName"`
	Expiration                 string `json:"expiration"`
	HashLength                 int    `json:"hashLength"`
	HashOffset                 int    `json:"hashOffset"`
	MatchAcrossPools           string `json:"matchAcrossPools"`
	MatchAcrossServices        string `json:"matchAcrossServices"`
	MatchAcrossVirtuals        string `json:"matchAcrossVirtuals"`
	Method                     string `json:"method"`
	Mirror                     string `json:"mirror"`
	Secure                     string `json:"secure"`
	TMPartition                string `json:"tmPartition"`
	OverrideConnectionLimit    string `json:"overrideConnectionLimit"`
	Timeout                    string `json:"timeout"`
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
