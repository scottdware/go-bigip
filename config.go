package bigip

import "encoding/json"

const (
	ACTIVE  = "ACTIVE"
	STANDBY = "STANDBY"

	INSYNC = "In Sync"
)

type Config struct {
	*BigIP
}

type Devices struct {
	Items []Device `json:"items"`
}

type Device struct {
	ActiveModules     []string `json:"activeModules,omitempty"`
	BaseMac           string   `json:"baseMac,omitempty"`
	Build             string   `json:"build,omitempty"`
	ChassisID         string   `json:"chassisId,omitempty"`
	ChassisType       string   `json:"chassisType,omitempty"`
	ConfigsyncIP      string   `json:"configsyncIp,omitempty"`
	Edition           string   `json:"edition,omitempty"`
	FailoverState     string   `json:"failoverState,omitempty"`
	FullPath          string   `json:"fullPath,omitempty"`
	Generation        int      `json:"generation,omitempty"`
	HaCapacity        int      `json:"haCapacity,omitempty"`
	Hostname          string   `json:"hostname,omitempty"`
	Kind              string   `json:"kind,omitempty"`
	ManagementIP      string   `json:"managementIp,omitempty"`
	MarketingName     string   `json:"marketingName,omitempty"`
	MirrorIP          string   `json:"mirrorIp,omitempty"`
	MirrorSecondaryIP string   `json:"mirrorSecondaryIp,omitempty"`
	MulticastIP       string   `json:"multicastIp,omitempty"`
	MulticastPort     int      `json:"multicastPort,omitempty"`
	Name              string   `json:"name,omitempty"`
	OptionalModules   []string `json:"optionalModules,omitempty"`
	Partition         string   `json:"partition,omitempty"`
	PlatformID        string   `json:"platformId,omitempty"`
	Product           string   `json:"product,omitempty"`
	SelfDevice        string   `json:"selfDevice,omitempty"`
	TimeZone          string   `json:"timeZone,omitempty"`
	UnicastAddress    []struct {
		EffectiveIP   string `json:"effectiveIp,omitempty"`
		EffectivePort int    `json:"effectivePort,omitempty"`
		IP            string `json:"ip,omitempty"`
		Port          int    `json:"port,omitempty"`
	} `json:"unicastAddress,omitempty"`
}

type TrafficGroups struct {
	Items []TrafficGroup `json:"items,omitempty"`
}

type TrafficGroup struct {
	AutoFailbackEnabled string   `json:"autoFailbackEnabled,omitempty"`
	AutoFailbackTime    int      `json:"autoFailbackTime,omitempty"`
	FullPath            string   `json:"fullPath,omitempty"`
	Generation          int      `json:"generation,omitempty"`
	HaLoadFactor        int      `json:"haLoadFactor,omitempty"`
	HAOrder             []string `json:"haOrder,omitempty"`
	IsFloating          string   `json:"isFloating,omitempty"`
	Kind                string   `json:"kind,omitempty"`
	Mac                 string   `json:"mac,omitempty"`
	Name                string   `json:"name,omitempty"`
	Partition           string   `json:"partition,omitempty"`
	UnitID              int      `json:"unitId,omitempty"`
}

// how am i supposed to parse that other than into a string
type failoverStatusEntries struct {
	Entries map[string]failoverStatus
}

type failoverStatus struct {
	NestedStats struct {
		Entries Entry `json:"entries"`
	} `json:"nestedStats"`
}

type syncStatusEntires struct {
	Entries map[string]syncStatus
}

type syncStatus struct {
	NestedStats struct {
		Entries Entry `json:"entries"`
	} `json:"nestedStats"`
}

type Entry struct {
	Color   ColorEntry   `json:"color,omitempty"`
	Mode    ModeEntry    `json:"mode,omitempty"`
	Status  StatusEntry  `json:"status,omitempty"`
	Summary SummaryEntry `json:"summary,omitempty"`
	Details DetailsEntry `json:"details,omitempty"`
}

type ColorEntry struct {
	Description string `json:"description"`
}

type ModeEntry struct {
	Description string `json:"description"`
}

type StatusEntry struct {
	Description string `json:"description"`
}

type SummaryEntry struct {
	Description string `json:"description"`
}

type DetailsEntry struct {
	Description string `json:"description"`
}

// assume one traffic group
// maybe a better idea is to parse cm/trafficGroup/stats
// accept a traffic group as a param?
// traffic group resources are stupid: cm/traffic-group/~Common~traffic-group-1:~Common~<hostname>/stats
func (c *Config) IsInSync() (bool, error) {
	ss, err := c.getSyncStatus()
	if err != nil {
		return false, err
	}
	for _, entry := range ss.Entries {
		if entry.NestedStats.Entries.Status.Description == INSYNC {
			return true, nil
		}
	}
	return false, nil
}

func (c *Config) getSyncStatus() (*syncStatusEntires, error) {
	var ss syncStatusEntires
	req := &APIRequest{
		Method: "get",
		URL:    "cm/syncStatus",
	}

	resp, err := c.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &ss)
	if err != nil {
		return nil, err
	}

	return &ss, nil
}

// assumes one traffic group
func (c *Config) IsActive() (bool, error) {
	fs, err := c.getFailoverStatus()
	if err != nil {
		return false, err
	}
	for _, entry := range fs.Entries {
		if entry.NestedStats.Entries.Status.Description == ACTIVE {
			return true, nil
		}
	}
	return false, nil
}

func (c *Config) getFailoverStatus() (*failoverStatusEntries, error) {
	var fs failoverStatusEntries
	req := &APIRequest{
		Method: "get",
		URL:    "cm/failoverStatus",
	}

	resp, err := c.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

func (c *Config) TrafficGroups() (*TrafficGroups, error) {
	var tg TrafficGroups
	req := &APIRequest{
		Method: "get",
		URL:    "cm/trafficGroup",
	}

	resp, err := c.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &tg)
	if err != nil {
		return nil, err
	}

	return &tg, nil
}

func (c *Config) Devices() (*Devices, error) {
	var d Devices
	req := &APIRequest{
		Method: "get",
		URL:    "cm/device",
	}

	resp, err := c.APICall(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
