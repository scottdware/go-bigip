package bigip

// IRules returns a list of irules
func (b *BigIP) IRules() (*IRules, error) {
	var rules IRules
	err, _ := b.getForEntity(&rules, uriLtm, uriIRule)
	if err != nil {
		return nil, err
	}

	return &rules, nil
}

// IRule returns information about the given iRule.
func (b *BigIP) IRule(name string) (*IRule, error) {
	var rule IRule
	err, ok := b.getForEntity(&rule, uriLtm, uriIRule, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &rule, nil
}

// CreateIRule creates a new iRule on the system.
func (b *BigIP) CreateIRule(name, rule string) error {
	irule := &IRule{
		Name: name,
		Rule: rule,
	}
	return b.post(irule, uriLtm, uriIRule)
}

// DeleteIRule removes an iRule from the system.
func (b *BigIP) DeleteIRule(name string) error {
	return b.delete(uriLtm, uriIRule, name)
}

// ModifyIRule updates the given iRule with any changed values.
func (b *BigIP) ModifyIRule(name string, irule *IRule) error {
	irule.Name = name
	return b.put(irule, uriLtm, uriIRule, name)
}

type IRules struct {
	IRules []IRule `json:"items"`
}

type IRule struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
	Rule      string `json:"apiAnonymous,omitempty"`
}
