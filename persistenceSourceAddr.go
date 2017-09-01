package bigip

type PersistenceSourceAddres struct {
	PersistenceSourceAddr []PersistenceSourceAddr `json:"items,omitempty"`
}

type PersistenceSourceAddr struct {
	Name                    string `json:"name"`
	DefaultsFrom            string `json:"defaultsFrom"`
	Kind                    string `json:"kind,omitempty"`
	Partition               string `json:"partition,omitempty"`
	FullPath                string `json:"fullPath,omitempty"`
	Generation              int    `json:"generation,omitempty"`
	SelfLink                string `json:"selfLink,omitempty"`
	HashAlgorithm           string `json:"hashAlgorithm,omitempty"`
	MapProxies              string `json:"mapProxies,omitempty"`
	MatchAcrossPools        string `json:"matchAcrossPools,omitempty"`
	MatchAcrossServices     string `json:"matchAcrossServices,omitempty"`
	MatchAcrossVirtuals     string `json:"matchAcrossVirtuals,omitempty"`
	Mirror                  string `json:"mirror,omitempty"`
	OverrideConnectionLimit string `json:"overrideConnectionLimit,omitempty"`
	Timeout                 string `json:"timeout,omitempty"`
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
