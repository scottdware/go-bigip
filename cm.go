package bigip

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Devices contains a list of every device on the BIG-IP system.
type Devices struct {
	Devices []Device `json:"items"`
}

// UnicastAddress represent a device unicast address
type UnicastAddress struct {
	EffectiveIP   string `json:"effectiveIp,omitempty"`
	EffectivePort int    `json:"effectivePort,omitempty"`
	IP            string `json:"ip,omitempty"`
	Port          int    `json:"port,omitempty"`
}

// Device contains information about each individual device.
type Device struct {
	Name               string           `json:"name,omitempty"`
	Partition          string           `json:"partition,omitempty"`
	FullPath           string           `json:"fullPath,omitempty"`
	Generation         int              `json:"generation,omitempty"`
	FailoverState      string           `json:"failoverState,omitempty"`
	Hostname           string           `json:"hostname,omitempty"`
	ManagementIp       string           `json:"managementIp,omitempty"`
	SelfDevice         string           `json:"selfDevice,omitempty"`
	ActiveModules      []string         `json:"activeModules,omitempty"`
	AppService         string           `json:"appService,omitempty"`
	BaseMac            string           `json:"baseMac,omitempty"`
	Build              string           `json:"build,omitempty"`
	Cert               string           `json:"cert,omitempty"`
	ChassisId          string           `json:"chassisId,omitempty"`
	ChassisType        string           `json:"chassisType,omitempty"`
	Comment            string           `json:"comment,omitempty"`
	ConfigSyncIP       string           `json:"configSyncIp,omitempty"`
	Contact            string           `json:"contact,omitempty"`
	Description        string           `json:"description,omitempty"`
	Edition            string           `json:"edition,omitempty"`
	HACapacity         int              `json:"haCapacity,omitempty"`
	InactiveModules    []string         `json:"inactiveModules,omitempty"`
	Key                string           `json:"key,omitempty"`
	Location           string           `json:"location,omitempty"`
	MarketingName      string           `json:"marketingName,omitempty"`
	MirrorIP           string           `json:"json:mirrorIp,omitempty"`
	MirrorSecondaryIP  string           `json:"mirrorSecondaryIp,omitempty"`
	MulticastInterface string           `json:"multicastInterface,omitempty"`
	MulticastIP        string           `json:"multicastIp,omitempty"`
	MulticastPort      int              `json:"multicastPort,omitempty"`
	OptionalModules    []string         `json:"optionalModules,omitempty"`
	TMPartition        string           `json:"tmPartition,omitempty"`
	PlatformID         string           `json:"platformId,omitempty"`
	Product            string           `json:"product,omitempty"`
	TimeLimitedModules []string         `json:"timeLimitedModules,omitempty"`
	Timezone           string           `json:"timeZone,omitempty"`
	UnicastAddress     []UnicastAddress `json:"unicastAddress,omitempty"`
	Version            string           `json:"version"`
}

type ConfigSync struct {
	Command     string `json:"command,omitempty"`
	UtilCmdArgs string `json:"utilCmdArgs,omitempty"`
}

const (
	uriCm                   = "cm"
	uriDevice               = "device"
	uriAutodeploy           = "autodeploy"
	uriSoftwareImageUploads = "software-image-uploads"
)

// Devices returns a list of devices.
func (b *BigIP) Devices() (*Devices, error) {
	var devices Devices
	err, _ := b.getForEntity(&devices, uriCm, uriDevice)

	if err != nil {
		return nil, err
	}

	return &devices, nil
}

// GetCurrentDevice returns a current device.
func (b *BigIP) GetCurrentDevice() (*Device, error) {
	devices, err := b.Devices()
	if err != nil {
		return nil, err
	}
	for _, d := range devices.Devices {
		// f5 api is returning bool value as string
		if d.SelfDevice == "true" {
			return &d, nil
		}
	}
	return nil, errors.New("could not find this device")
}

// ConfigSyncToGroup runs command config-sync to-group <attr>
func (b *BigIP) ConfigSyncToGroup(name string) error {
	args := "config-sync to-group " + name
	config := &ConfigSync{
		Command:     "run",
		UtilCmdArgs: args,
	}
	return b.post(config, uriCm)
}

// Upload a software image
func (b *BigIP) UploadSoftwareImage(f *os.File) (*Upload, error) {
	if !strings.HasSuffix(f.Name(), ".iso") {
		err := fmt.Errorf("File must have .iso extension")
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return b.Upload(f, info.Size(), uriCm, uriAutodeploy, uriSoftwareImageUploads, info.Name())
}
