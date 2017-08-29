package bigip

import "encoding/json"

// DataGroups contains a list of data groups on the BIG-IP system.
type DataGroups struct {
	SelfLink   string      `json:"selfLink,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	DataGroups []DataGroup `json:"items,omitempty"`
}

// DataGroups contains information about each data group.
type DataGroup struct {
	Kind       string            `json:"kind,omitempty"`
	Name       string            `json:"name"`
	FullPath   string            `json:"fullPath,omitempty"`
	Partition  string            `json:"tmPartition,omitempty"`
	Generation int               `json:"generation,omitempty"`
	SelfLink   string            `json:"selfLink,omitempty"`
	Type       string            `json:"type,omitempty"`
	Records    []DataGroupRecord `json:"records"`
}

type DataGroupRecord struct {
	Name      string `json:"name"`
	Partition string `json:"partition,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
	Data      string `json:"data"`
}

type dataGroupDTO struct {
	Kind       string            `json:"kind,omitempty"`
	Name       string            `json:"name,omitempty"`
	FullPath   string            `json:"fullPath,omitempty"`
	Partition  string            `json:"tmPartition,omitempty"`
	Generation int               `json:"generation,omitempty"`
	SelfLink   string            `json:"selfLink,omitempty"`
	Type       string            `json:"type,omitempty"`
	Records    []DataGroupRecord `json:"records,omitempty"`
}

func (p *DataGroup) MarshalJSON() ([]byte, error) {
	var dto dataGroupDTO
	marshal(&dto, p)
	return json.Marshal(dto)
}

func (p *DataGroup) UnmarshalJSON(b []byte) error {
	var dto dataGroupDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(p, &dto)
}
