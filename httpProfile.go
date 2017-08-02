package bigip

type HttpProfiles struct {
	HttpProfiles []HttpProfile `json:"items"`
}

type HttpProfile struct {
	Kind         string `json:"kind"`
	DefaultsFrom string `json:"defaultsFrom"`
	Name         string `json:"name"`
	Partition    string `json:"partition"`
	FullPath     string `json:"fullPath"`
	Generation   int    `json:"generation"`
	SelfLink     string `json:"selfLink"`
	AcceptXff    string `json:"acceptXff"`
	Enforcement  struct {
		ExcessClientHeaders   string `json:"excessClientHeaders"`
		ExcessServerHeaders   string `json:"excessServerHeaders"`
		MaxHeaderCount        int    `json:"maxHeaderCount"`
		MaxHeaderSize         int    `json:"maxHeaderSize"`
		MaxRequests           int    `json:"maxRequests"`
		OversizeClientHeaders string `json:"oversizeClientHeaders"`
		OversizeServerHeaders string `json:"oversizeServerHeaders"`
		Pipeline              string `json:"pipeline"`
		TruncatedRedirects    string `json:"truncatedRedirects"`
		UnknownMethod         string `json:"unknownMethod"`
	} `json:"enforcement"`
	ExplicitProxy struct {
		DefaultConnectHandling string `json:"defaultConnectHandling"`
	} `json:"explicitProxy"`
	InsertXforwardedFor       string `json:"insertXforwardedFor"`
	LwsWidth                  int    `json:"lwsWidth"`
	OneconnectTransformations string `json:"oneconnectTransformations"`
	ProxyType                 string `json:"proxyType"`
	RequestChunking           string `json:"requestChunking"`
	ResponseChunking          string `json:"responseChunking"`
	ServerAgentName           string `json:"serverAgentName"`
	Sflow                     struct {
		PollInterval       int    `json:"pollInterval"`
		PollIntervalGlobal string `json:"pollIntervalGlobal"`
		SamplingRate       int    `json:"samplingRate"`
		SamplingRateGlobal string `json:"samplingRateGlobal"`
	} `json:"sflow"`
	ViaRequest  string `json:"viaRequest"`
	ViaResponse string `json:"viaResponse"`
}

// HttpProfiles returns a list of http profiles.
func (b *BigIP) HttpProfiles() (*HttpProfiles, error) {
	var httpProfiles HttpProfiles
	err, _ := b.getForEntity(&httpProfiles, uriLtm, uriProfile, uriProfileHttp)
	if err != nil {
		return nil, err
	}

	return &httpProfiles, nil
}

// GetHttpProfile gets a http profile by name. Returns nil if the http profile does not exist
func (b *BigIP) GetHttpProfile(name string) (*HttpProfile, error) {
	var httpProfile HttpProfile
	err, ok := b.getForEntity(&httpProfile, uriLtm, uriProfile, uriProfileHttp, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &httpProfile, nil
}

// CreateHttpProfile creates a new http profile on the BIG-IP system.
func (b *BigIP) CreateHttpProfile(name string, parent string) error {
	config := &HttpProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriProfileHttp)
}

// AddHttpProfile adds a new http profile on the BIG-IP system.
func (b *BigIP) AddHttpProfile(config *HttpProfile) error {
	return b.post(config, uriLtm, uriProfile, uriProfileHttp)
}

// DeleteHttpProfile removes a http profile.
func (b *BigIP) DeleteHttpProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriProfileHttp, name)
}

// ModifyHttpProfile allows you to change any attribute of a http profile.
// Fields that can be modified are referenced in the HttpProfile struct.
func (b *BigIP) ModifyHttpProfile(name string, config *HttpProfile) error {
	return b.put(config, uriLtm, uriProfile, uriProfileHttp, name)
}
