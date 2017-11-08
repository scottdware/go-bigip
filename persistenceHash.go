package bigip

type PersistenceHashes struct {
	PersistenceHash []PersistenceHash `json:"items,omitempty"`
}

type PersistenceHash struct {
	Kind                    string `json:"kind,omitempty"`
	DefaultsFrom            string `json:"defaultsFrom"`
	Name                    string `json:"name"`
	Partition               string `json:"partition,omitempty"`
	FullPath                string `json:"fullPath,omitempty"`
	Generation              int    `json:"generation,omitempty"`
	SelfLink                string `json:"selfLink,omitempty"`
	HashAlgorithm           string `json:"hashAlgorithm,omitempty"`
	HashBufferLimit         int    `json:"hashBufferLimit,omitempty"`
	HashLength              int    `json:"hashLength,omitempty"`
	HashOffset              int    `json:"hashOffset,omitempty"`
	MatchAcrossPools        string `json:"matchAcrossPools,omitempty"`
	MatchAcrossServices     string `json:"matchAcrossServices,omitempty"`
	MatchAcrossVirtuals     string `json:"matchAcrossVirtuals,omitempty"`
	AppService              string `json:"appService,omitempty"`
	Description             string `json:"description,omitempty"`
	HashEndPattern          string `json:"hashEndPattern,omitempty"`
	HashStartPattern        string `json:"hashStartPattern,omitempty"`
	Mode                    string `json:"mode,omitempty"`
	Rule                    string `json:"rule,omitempty"`
	TMPartition             string `json:"tmPartition,omitempty"`
	Mirror                  string `json:"mirror,omitempty"`
	OverrideConnectionLimit string `json:"overrideConnectionLimit,omitempty"`
	Timeout                 string `json:"timeout,omitempty"`
}

// PersistenceHash returns a list of oersistence profiles.
func (b *BigIP) PersistenceHash() (*PersistenceHash, error) {
	var persistenceProfiles PersistenceHash
	err, _ := b.getForEntity(&persistenceProfiles, uriLtm, uriPersistences)
	if err != nil {
		return nil, err
	}

	return &persistenceProfiles, nil
}

// GetPersistenceHash gets a persistence profile by name. Returns nil if the persistence profile does not exist
func (b *BigIP) GetPersistenceHash(name string) (*PersistenceHash, error) {
	var persistenceProfile PersistenceHash
	err, ok := b.getForEntity(&persistenceProfile, uriLtm, uriProfile, uriPersistenceHash, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &persistenceProfile, nil
}

// CreatePersistenceHash creates a new persistence profile on the BIG-IP system.
func (b *BigIP) CreatePersistenceHash(name string, parent string) error {
	config := &PersistenceHash{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriPersistenceHash)
}

// AddPersistenceHash adds a new persistence profile on the BIG-IP system.
func (b *BigIP) AddPersistenceHash(config *PersistenceHash) error {
	return b.post(config, uriLtm, uriProfile, uriPersistenceHash)
}

// DeletePersistenceHash removes a persistence profile.
func (b *BigIP) DeletePersistenceHash(name string) error {
	return b.delete(uriLtm, uriProfile, uriPersistenceHash, name)
}

// ModifyPersistenceHash allows you to change any attribute of a persistence profile.
// Fields that can be modified are referenced in the PersistenceHash struct.
func (b *BigIP) ModifyPersistenceHash(name string, config *PersistenceHash) error {
	return b.put(config, uriLtm, uriProfile, uriPersistenceHash, name)
}
