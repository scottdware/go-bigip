package bigip

import "encoding/json"

const (
	uriSys            = "sys"
	uriFolder         = "folder"
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

// Folders contains a list of every folder on the BIG-IP system.
type Folders struct {
	Folders []Folder `json:"items"`
}

type folderDTO struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`

	AppService  string `json:"appService,omitempty"`
	Description string `json:"description,omitempty"`
	// Set to "default" to inherit or a device group name to control. You can also set it to "non-default" to pin its device group to its current setting and turn off inheritance.
	DeviceGroup string `json:"deviceGroup,omitempty"`
	Hidden      string `json:"hidden,omitempty" bool:"true"`
	NoRefCheck  string `json:"noRefCheck,omitempty" bool:"true"`
	// Set to "default" to inherit or a traffic group name to control. You can also set it to "non-default" to pin its traffic group to its current setting and turn off inheritance.
	TrafficGroup string `json:"trafficGroup,omitempty"`

	// Read-only property. Set DeviceGroup to control.
	InheritedDeviceGroup string `json:"inheritedDevicegroup,omitempty" bool:"true"`

	// Read-only property. Set TrafficGroup to control.
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty" bool:"true"`
}

type Folder struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`

	AppService   string `json:"appService,omitempty"`
	Description  string `json:"description,omitempty"`
	DeviceGroup  string `json:"deviceGroup,omitempty"`
	Hidden       *bool  `json:"hidden,omitempty"`
	NoRefCheck   *bool  `json:"noRefCheck,omitempty"`
	TrafficGroup string `json:"trafficGroup,omitempty"`

	// Read-only property. Set DeviceGroup to "default" or "non-default" to control.
	InheritedDeviceGroup *bool `json:"inheritedDevicegroup,omitempty"`

	// Read-only property. Set TrafficGroup to "default" or "non-default" to control.
	InheritedTrafficGroup *bool `json:"inheritedTrafficGroup,omitempty"`
}

func (f *Folder) MarshalJSON() ([]byte, error) {
	var dto folderDTO
	marshal(&dto, f)
	return json.Marshal(dto)
}

func (f *Folder) UnmarshalJSON(b []byte) error {
	var dto folderDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(f, &dto)
}

// Folders returns a list of folders.
func (b *BigIP) Folders() (*Folders, error) {
	var folders Folders
	err, _ := b.getForEntity(&folders, uriSys, uriFolder)
	if err != nil {
		return nil, err
	}

	return &folders, nil
}

// CreateFolder adds a new folder to the BIG-IP system.
func (b *BigIP) CreateFolder(name string) error {
	config := &Folder{
		Name: name,
	}

	return b.post(config, uriSys, uriFolder)
}

// AddFolder adds a new folder by config to the BIG-IP system.
func (b *BigIP) AddFolder(config *Folder) error {

	return b.post(config, uriSys, uriFolder)
}

// GetFolder retrieves a Folder by name. Returns nil if the folder does not exist
func (b *BigIP) GetFolder(name string) (*Folder, error) {
	var folder Folder
	err, ok := b.getForEntity(&folder, uriSys, uriFolder, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &folder, nil
}

// DeleteFolder removes a folder.
func (b *BigIP) DeleteFolder(name string) error {
	return b.delete(uriSys, uriFolder, name)
}

// ModifyFolder allows you to change any attribute of a folder. Fields that can
// be modified are referenced in the Folder struct. This replaces the existing
// configuration, so use PatchFolder if you want to change only particular
// attributes.
func (b *BigIP) ModifyFolder(name string, config *Folder) error {
	return b.put(config, uriSys, uriFolder, name)
}

// PatchFolder allows you to change any attribute of a folder. Fields that can
// be modified are referenced in the Folder struct. This changes only the
// attributes provided, so use ModifyFolder if you want to replace the existing
// configuration.
func (b *BigIP) PatchFolder(name string, config *Folder) error {
	return b.patch(config, uriSys, uriFolder, name)
}
