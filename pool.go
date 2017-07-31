package bigip

// Pool contains information about each pool. You can use all of these
// fields when modifying a pool.
type Pool struct {
	Name                   string `json:"name,omitempty"`
	Partition              string `json:"partition,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	AllowNAT               string `json:"allowNat,omitempty"`
	AllowSNAT              string `json:"allowSnat,omitempty"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight,omitempty"`
	IPTOSToClient          string `json:"ipTosToClient,omitempty"`
	IPTOSToServer          string `json:"ipTosToServer,omitempty"`
	LinkQoSToClient        string `json:"linkQosToClient,omitempty"`
	LinkQoSToServer        string `json:"linkQosToServer,omitempty"`
	LoadBalancingMode      string `json:"loadBalancingMode,omitempty"`
	MinActiveMembers       int    `json:"minActiveMembers,omitempty"`
	MinUpMembers           int    `json:"minUpMembers,omitempty"`
	MinUpMembersAction     string `json:"minUpMembersAction,omitempty"`
	MinUpMembersChecking   string `json:"minUpMembersChecking,omitempty"`
	Monitor                string `json:"monitor,omitempty"`
	QueueDepthLimit        int    `json:"queueDepthLimit,omitempty"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit,omitempty"`
	QueueTimeLimit         int    `json:"queueTimeLimit,omitempty"`
	ReselectTries          int    `json:"reselectTries,omitempty"`
	ServiceDownAction      string `json:"serviceDownAction,omitempty"`
	SlowRampTime           int    `json:"slowRampTime,omitempty"`
}

// Pool Members contains a list of pool members within a pool on the BIG-IP system.
type PoolMembers struct {
	PoolMembers []PoolMember `json:"items"`
}

// poolMember is used only when adding members to a pool.
type poolMember struct {
	Name string `json:"name"`
}

// poolMembers is used only when modifying members on a pool.
type poolMembers struct {
	Members []PoolMember `json:"members"`
}

// Pool Member contains information about each individual member in a pool. You can use all
// of these fields when modifying a pool member.
type PoolMember struct {
	Name            string `json:"name,omitempty"`
	Partition       string `json:"partition,omitempty"`
	FullPath        string `json:"fullPath,omitempty"`
	Generation      int    `json:"generation,omitempty"`
	Address         string `json:"address,omitempty"`
	ConnectionLimit int    `json:"connectionLimit,omitempty"`
	DynamicRatio    int    `json:"dynamicRatio,omitempty"`
	InheritProfile  string `json:"inheritProfile,omitempty"`
	Logging         string `json:"logging,omitempty"`
	Monitor         string `json:"monitor,omitempty"`
	PriorityGroup   int    `json:"priorityGroup,omitempty"`
	RateLimit       string `json:"rateLimit,omitempty"`
	Ratio           int    `json:"ratio,omitempty"`
	Session         string `json:"session,omitempty"`
	State           string `json:"state,omitempty"`
}

// PoolMembers returns a list of pool members for the given pool.
func (b *BigIP) PoolMembers(name string) (*PoolMembers, error) {
	var poolMembers PoolMembers
	err, _ := b.getForEntity(&poolMembers, uriLtm, uriPool, name, uriPoolMember)
	if err != nil {
		return nil, err
	}

	return &poolMembers, nil
}

// AddPoolMember adds a node/member to the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) AddPoolMember(pool, member string) error {
	config := &poolMember{
		Name: member,
	}

	return b.post(config, uriLtm, uriPool, pool, uriPoolMember)
}

// GetPoolMember returns the details of a member in the specified pool.
func (b *BigIP) GetPoolMember(pool string, member string) (*PoolMember, error) {
	var poolMember PoolMember
	err, ok := b.getForEntity(&poolMember, uriLtm, uriPool, pool, uriPoolMember, member)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &poolMember, nil
}

// CreatePoolMember creates a pool member for the specified pool.
func (b *BigIP) CreatePoolMember(pool string, config *PoolMember) error {
	return b.post(config, uriLtm, uriPool, pool, uriPoolMember)
}

// ModifyPoolMember will update the configuration of a particular pool member.
func (b *BigIP) ModifyPoolMember(pool string, config *PoolMember) error {
	member := config.FullPath
	// These fields are not used when modifying a pool member; so omit them.
	config.Name = ""
	config.Partition = ""
	config.FullPath = ""

	// This cannot be modified for an existing pool member.
	config.Address = ""

	return b.put(config, uriLtm, uriPool, pool, uriPoolMember, member)
}

// UpdatePoolMembers does a replace-all-with for the members of a pool.
func (b *BigIP) UpdatePoolMembers(pool string, pm *[]PoolMember) error {
	config := &poolMembers{
		Members: *pm,
	}
	return b.put(config, uriLtm, uriPool, pool)
}

// RemovePoolMember removes a pool member from the specified pool.
func (b *BigIP) RemovePoolMember(pool string, config *PoolMember) error {
	member := config.FullPath
	return b.delete(uriLtm, uriPool, pool, uriPoolMember, member)
}

// DeletePoolMember removes a member from the given pool. <member> must be in the form
// of <node>:<port>, i.e.: "web-server1:443".
func (b *BigIP) DeletePoolMember(pool string, member string) error {
	return b.delete(uriLtm, uriPool, pool, uriPoolMember, member)
}

// PoolMemberStatus changes the status of a pool member. <state> can be either
// "enable" or "disable". <member> must be in the form of <node>:<port>,
// i.e.: "web-server1:443".
func (b *BigIP) PoolMemberStatus(pool string, member string, state string) error {
	config := &Node{}

	switch state {
	case "enable":
		// config.State = "unchecked"
		config.Session = "user-enabled"
	case "disable":
		// config.State = "unchecked"
		config.Session = "user-disabled"
		// case "offline":
		// 	config.State = "user-down"
		// 	config.Session = "user-disabled"
	}

	return b.put(config, uriLtm, uriPool, pool, uriPoolMember, member)
}

// CreatePool adds a new pool to the BIG-IP system by name.
func (b *BigIP) CreatePool(name string) error {
	config := &Pool{
		Name: name,
	}

	return b.post(config, uriLtm, uriPool)
}

// AddPool creates a new pool on the BIG-IP system.
func (b *BigIP) AddPool(config *Pool) error {
	return b.post(config, uriLtm, uriPool)
}

// Get a Pool by name. Returns nil if the Pool does not exist
func (b *BigIP) GetPool(name string) (*Pool, error) {
	var pool Pool
	err, ok := b.getForEntity(&pool, uriLtm, uriPool, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &pool, nil
}

// DeletePool removes a pool.
func (b *BigIP) DeletePool(name string) error {
	return b.delete(uriLtm, uriPool, name)
}

// ModifyPool allows you to change any attribute of a pool. Fields that
// can be modified are referenced in the Pool struct.
func (b *BigIP) ModifyPool(name string, config *Pool) error {
	return b.put(config, uriLtm, uriPool, name)
}

// Pools returns a list of pools.
func (b *BigIP) Pools() (*Pools, error) {
	var pools Pools
	err, _ := b.getForEntity(&pools, uriLtm, uriPool)
	if err != nil {
		return nil, err
	}

	return &pools, nil
}
