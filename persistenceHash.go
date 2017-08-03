package bigip

type PersistenceHashes struct {
	PersistenceHash []PersistenceHash `json:"items"`
}

type PersistenceHash struct {
	Kind                    string `json:"kind"`
	DefaultsFrom            string `json:"defaultsFrom"`
	Name                    string `json:"name"`
	Partition               string `json:"partition"`
	FullPath                string `json:"fullPath"`
	Generation              int    `json:"generation"`
	SelfLink                string `json:"selfLink"`
	HashAlgorithm           string `json:"hashAlgorithm"`
	HashBufferLimit         int    `json:"hashBufferLimit"`
	HashLength              int    `json:"hashLength"`
	HashOffset              int    `json:"hashOffset"`
	MatchAcrossPools        string `json:"matchAcrossPools"`
	MatchAcrossServices     string `json:"matchAcrossServices"`
	MatchAcrossVirtuals     string `json:"matchAcrossVirtuals"`
	AppService              string `json:"appService"`
	Description             string `json:"description"`
	HashEndPattern          string `json:"hashEndPattern"`
	HashStartPattern        string `json:"hashStartPattern"`
	Mode                    string `json:"mode"`
	Rule                    string `json:"rule"`
	TMPartition             string `json:"tmPartition"`
	Mirror                  string `json:"mirror"`
	OverrideConnectionLimit string `json:"overrideConnectionLimit"`
	Timeout                 string `json:"timeout"`
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
