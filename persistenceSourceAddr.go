package bigip

type PersistenceSourceAddres struct {
	PersistenceSourceAddr []PersistenceSourceAddr `json:"items"`
}

type PersistenceSourceAddr struct {
	Name                    string `json:"name"`
	DefaultsFrom            string `json:"defaultsFrom"`
	Kind                    string `json:"kind"`
	Partition               string `json:"partition"`
	FullPath                string `json:"fullPath"`
	Generation              int    `json:"generation"`
	SelfLink                string `json:"selfLink"`
	HashAlgorithm           string `json:"hashAlgorithm"`
	MapProxies              string `json:"mapProxies"`
	MatchAcrossPools        string `json:"matchAcrossPools"`
	MatchAcrossServices     string `json:"matchAcrossServices"`
	MatchAcrossVirtuals     string `json:"matchAcrossVirtuals"`
	Mirror                  string `json:"mirror"`
	OverrideConnectionLimit string `json:"overrideConnectionLimit"`
	Timeout                 string `json:"timeout"`
}

// PersistenceSourceAddr returns a list of oersistence profiles.
func (b *BigIP) PersistenceSourceAddr() (*PersistenceSourceAddr, error) {
	var persistenceProfiles PersistenceSourceAddr
	err, _ := b.getForEntity(&persistenceProfiles, uriLtm, uriPersistences)
	if err != nil {
		return nil, err
	}

	return &persistenceProfiles, nil
}

// GetPersistenceSourceAddr gets a persistence profile by name. Returns nil if the persistence profile does not exist
func (b *BigIP) GetPersistenceSourceAddr(name string) (*PersistenceSourceAddr, error) {
	var persistenceProfile PersistenceSourceAddr
	err, ok := b.getForEntity(&persistenceProfile, uriLtm, uriProfile, uriPersistenceSourceAddr, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &persistenceProfile, nil
}

// CreatePersistenceSourceAddr creates a new persistence profile on the BIG-IP system.
func (b *BigIP) CreatePersistenceSourceAddr(name string, parent string) error {
	config := &PersistenceSourceAddr{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriPersistenceSourceAddr)
}

// AddPersistenceSourceAddr adds a new persistence profile on the BIG-IP system.
func (b *BigIP) AddPersistenceSourceAddr(config *PersistenceSourceAddr) error {
	return b.post(config, uriLtm, uriProfile, uriPersistenceSourceAddr)
}

// DeletePersistenceSourceAddr removes a persistence profile.
func (b *BigIP) DeletePersistenceSourceAddr(name string) error {
	return b.delete(uriLtm, uriProfile, uriPersistenceSourceAddr, name)
}

// ModifyPersistenceSourceAddr allows you to change any attribute of a persistence profile.
// Fields that can be modified are referenced in the PersistenceSourceAddr struct.
func (b *BigIP) ModifyPersistenceSourceAddr(name string, config *PersistenceSourceAddr) error {
	return b.put(config, uriLtm, uriProfile, uriPersistenceSourceAddr, name)
}
