package bigip

import (
	"errors"
)

// Devices contains a list of every device on the BIG-IP system.
type Devices struct {
	Devices []Device `json:"items"`
}


// Device contains information about each individual device.
type Device struct {
	Name          string `json:"name,omitempty"`
	Partition     string `json:"partition,omitempty"`
	FullPath      string `json:"fullPath,omitempty"`
	Generation    int    `json:"generation,omitempty"`
	FailoverState string `json:"failoverState,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	ManagementIp  string `json:"managementIp,omitempty"`
	SelfDevice    string `json:"selfDevice,omitempty"`
}

type ConfigSync struct {
	Command     string `json:"command,omitempty"`
	UtilCmdArgs string `json:"utilCmdArgs,omitempty"`
}

const (
	uriCm     = "cm"
	uriDevice = "device"
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
	args := "config-sync to-group "+name
	config := &ConfigSync{
		Command: "run",
		UtilCmdArgs: args,
	}
	return b.post(config, uriCm)
}
