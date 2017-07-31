package bigip

import "encoding/json"

// DataGroups contains a list of data groups on the BIG-IP system.
type DataGroups struct {
	DataGroups []DataGroup `json:"items"`
}

// DataGroups contains information about each data group.
type DataGroup struct {
	Name       string
	Partition  string
	FullPath   string
	Generation int
	Type       string
	Records    []DataGroupRecord
}

type DataGroupRecord struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}

type dataGroupDTO struct {
	Name       string            `json:"name,omitempty"`
	Partition  string            `json:"partition,omitempty"`
	FullPath   string            `json:"fullPath,omitempty"`
	Generation int               `json:"generation,omitempty"`
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
