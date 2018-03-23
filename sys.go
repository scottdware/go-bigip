package bigip

const (
	uriSys            = "sys"
	uriSyslog         = "syslog"
	uriSoftware       = "software"
	uriVolume         = "volume"
	uriHardware       = "hardware"
	uriGlobalSettings = "global-settings"
	uriManagementIp   = "management-ip"
	//uriPlatform = "?$select=platform"
)

type Volumes struct {
	Volumes []Volume `json:"items,omitempty"`
}

type Volume struct {
	Name       string `json:"items,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	SelfLink   string `json:"selfLink,omitempty"`
	Active     bool   `json:"active,omitempty"`
	BaseBuild  string `json:"basebuild,omitempty"`
	Build      string `json:"build,omitempty"`
	Product    string `json:"product,omitempty"`
	Status     string `json:"status,omitempty"`
	Version    string `json:"version,omitempty"`
}

// Volumes returns a list of Software Volumes.
func (b *BigIP) Volumes() (*Volumes, error) {
	var volumes Volumes
	err, _ := b.getForEntity(&volumes, uriSys, uriSoftware, uriVolume)
	if err != nil {
		return nil, err
	}

	return &volumes, nil
}

type ManagementIP struct {
	Addresses []ManagementIPAddress
}

type ManagementIPAddress struct {
	Name       string `json:"items,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	SelfLink   string `json:"selfLink,omitempty"`
}

func (b *BigIP) ManagementIPs() (*ManagementIP, error) {
	var managementIP ManagementIP
	err, _ := b.getForEntity(&managementIP, uriSys, uriManagementIp)
	if err != nil {
		return nil, err
	}

	return &managementIP, nil
}

type SyslogRemoteServer struct {
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	LocalIP    string `json:"localIp,omitempty"`
	RemotePort int    `json:"remotePort,omitempty"`
}

type Syslog struct {
	SelfLink      string               `json:"selfLink,omitempty"`
	RemoteServers []SyslogRemoteServer `json:"remoteServers,omitempty"`
}

func (b *BigIP) Syslog() (*Syslog, error) {
	var syslog Syslog

	err, _ := b.getForEntity(&syslog, uriSys, uriSyslog)
	if err != nil {
		return nil, err
	}

	return &syslog, nil
}

func (b *BigIP) SetSyslog(config Syslog) error {
	return b.put(config, uriSys, uriSyslog)
}
