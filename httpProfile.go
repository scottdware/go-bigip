package bigip

type HttpProfiles struct {
	HttpProfiles []HttpProfile `json:"items"`
}

type HttpProfile struct {
	Kind         string `json:"kind,omitempty"`
	DefaultsFrom string `json:"defaultsFrom"`
	Name         string `json:"name"`
	Partition    string `json:"partition,omitempty"`
	FullPath     string `json:"fullPath,omitempty"`
	Generation   int    `json:"generation,omitempty"`
	SelfLink     string `json:"selfLink,omitempty"`
	AcceptXff    string `json:"acceptXff,omitempty"`
	Enforcement  struct {
		ExcessClientHeaders   string `json:"excessClientHeaders,omitempty"`
		ExcessServerHeaders   string `json:"excessServerHeaders,omitempty"`
		MaxHeaderCount        int    `json:"maxHeaderCount,omitempty"`
		MaxHeaderSize         int    `json:"maxHeaderSize,omitempty"`
		MaxRequests           int    `json:"maxRequests,omitempty"`
		OversizeClientHeaders string `json:"oversizeClientHeaders,omitempty"`
		OversizeServerHeaders string `json:"oversizeServerHeaders,omitempty"`
		Pipeline              string `json:"pipeline,omitempty"`
		TruncatedRedirects    string `json:"truncatedRedirects,omitempty"`
		UnknownMethod         string `json:"unknownMethod,omitempty"`
	} `json:"enforcement,omitempty"`
	ExplicitProxy struct {
		DefaultConnectHandling string `json:"defaultConnectHandling,omitempty"`
	} `json:"explicitProxy,omitempty"`
	InsertXforwardedFor       string `json:"insertXforwardedFor,omitempty"`
	LwsWidth                  int    `json:"lwsWidth,omitempty"`
	OneconnectTransformations string `json:"oneconnectTransformations,omitempty"`
	ProxyType                 string `json:"proxyType,omitempty"`
	RequestChunking           string `json:"requestChunking,omitempty"`
	ResponseChunking          string `json:"responseChunking,omitempty"`
	ServerAgentName           string `json:"serverAgentName,omitempty"`
	Sflow                     struct {
		PollInterval       int    `json:"pollInterval,omitempty"`
		PollIntervalGlobal string `json:"pollIntervalGlobal,omitempty"`
		SamplingRate       int    `json:"samplingRate,omitempty"`
		SamplingRateGlobal string `json:"samplingRateGlobal,omitempty"`
	} `json:"sflow,omitempty"`
	ViaRequest  string `json:"viaRequest,omitempty"`
	ViaResponse string `json:"viaResponse,omitempty"`
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
