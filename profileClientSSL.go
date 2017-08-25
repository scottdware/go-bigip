package bigip

// ClientSSLProfiles
// Documentation: https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_ltm_profile_client-ssl.ashx

// ClientSSLProfiles contains a list of every client-ssl profile on the BIG-IP system.
type ClientSSLProfiles struct {
	ClientSSLProfiles []ClientSSLProfile `json:"items"`
}

// ClientSSLProfile contains information about each client-ssl profile. You can use all
// of these fields when modifying a client-ssl profile.
type ClientSSLProfile struct {
	Name              string `json:"name,omitempty"`
	Partition         string `json:"partition,omitempty"`
	FullPath          string `json:"fullPath,omitempty"`
	Generation        int    `json:"generation,omitempty"`
	AlertTimeout      string `json:"alertTimeout,omitempty"`
	AllowNonSsl       string `json:"allowNonSsl,omitempty"`
	Authenticate      string `json:"authenticate,omitempty"`
	AuthenticateDepth int    `json:"authenticateDepth,omitempty"`
	CaFile            string `json:"caFile,omitempty"`
	CacheSize         int    `json:"cacheSize,omitempty"`
	CacheTimeout      int    `json:"cacheTimeout,omitempty"`
	Cert              string `json:"cert,omitempty"`
	CertKeyChain      []struct {
		Name       string `json:"name,omitempty"`
		Cert       string `json:"cert,omitempty"`
		Chain      string `json:"chain,omitempty"`
		Key        string `json:"key,omitempty"`
		Passphrase string `json:"passphrase,omitempty"`
	} `json:"certKeyChain,omitempty"`
	CertExtensionIncludes           []string `json:"certExtensionIncludes,omitempty"`
	CertLifespan                    int      `json:"certLifespan,omitempty"`
	CertLookupByIpaddrPort          string   `json:"certLookupByIpaddrPort,omitempty"`
	Chain                           string   `json:"chain,omitempty"`
	Ciphers                         string   `json:"ciphers,omitempty"`
	ClientCertCa                    string   `json:"clientCertCa,omitempty"`
	CrlFile                         string   `json:"crlFile,omitempty"`
	DefaultsFrom                    string   `json:"defaultsFrom,omitempty"`
	ForwardProxyBypassDefaultAction string   `json:"forwardProxyBypassDefaultAction,omitempty"`
	GenericAlert                    string   `json:"genericAlert,omitempty"`
	HandshakeTimeout                string   `json:"handshakeTimeout,omitempty"`
	InheritCertkeychain             string   `json:"inheritCertkeychain,omitempty"`
	Key                             string   `json:"key,omitempty"`
	ModSslMethods                   string   `json:"modSslMethods,omitempty"`
	Mode                            string   `json:"mode,omitempty"`
	TmOptions                       []string `json:"tmOptions,omitempty"`
	Passphrase                      string   `json:"passphrase,omitempty"`
	PeerCertMode                    string   `json:"peerCertMode,omitempty"`
	ProxyCaCert                     string   `json:"proxyCaCert,omitempty"`
	ProxyCaKey                      string   `json:"proxyCaKey,omitempty"`
	ProxyCaPassphrase               string   `json:"proxyCaPassphrase,omitempty"`
	ProxySsl                        string   `json:"proxySsl,omitempty"`
	ProxySslPassthrough             string   `json:"proxySslPassthrough,omitempty"`
	RenegotiatePeriod               string   `json:"renegotiatePeriod,omitempty"`
	RenegotiateSize                 string   `json:"renegotiateSize,omitempty"`
	Renegotiation                   string   `json:"renegotiation,omitempty"`
	RetainCertificate               string   `json:"retainCertificate,omitempty"`
	SecureRenegotiation             string   `json:"secureRenegotiation,omitempty"`
	ServerName                      string   `json:"serverName,omitempty"`
	SessionMirroring                string   `json:"sessionMirroring,omitempty"`
	SessionTicket                   string   `json:"sessionTicket,omitempty"`
	SniDefault                      string   `json:"sniDefault,omitempty"`
	SniRequire                      string   `json:"sniRequire,omitempty"`
	SslForwardProxy                 string   `json:"sslForwardProxy,omitempty"`
	SslForwardProxyBypass           string   `json:"sslForwardProxyBypass,omitempty"`
	SslSignHash                     string   `json:"sslSignHash,omitempty"`
	StrictResume                    string   `json:"strictResume,omitempty"`
	UncleanShutdown                 string   `json:"uncleanShutdown,omitempty"`
}

// ClientSSLProfiles returns a list of client-ssl profiles.
func (b *BigIP) ClientSSLProfiles() (*ClientSSLProfiles, error) {
	var clientSSLProfiles ClientSSLProfiles
	err, _ := b.getForEntity(&clientSSLProfiles, uriLtm, uriProfile, uriClientSSL)
	if err != nil {
		return nil, err
	}

	return &clientSSLProfiles, nil
}

// GetClientSSLProfile gets a client-ssl profile by name. Returns nil if the client-ssl profile does not exist
func (b *BigIP) GetClientSSLProfile(name string) (*ClientSSLProfile, error) {
	var clientSSLProfile ClientSSLProfile
	err, ok := b.getForEntity(&clientSSLProfile, uriLtm, uriProfile, uriClientSSL, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &clientSSLProfile, nil
}

// CreateClientSSLProfile creates a new client-ssl profile on the BIG-IP system.
func (b *BigIP) CreateClientSSLProfile(name string, parent string) error {
	config := &ClientSSLProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriClientSSL)
}

// AddClientSSLProfile adds a new client-ssl profile on the BIG-IP system.
func (b *BigIP) AddClientSSLProfile(config *ClientSSLProfile) error {
	return b.post(config, uriLtm, uriProfile, uriClientSSL)
}

// DeleteClientSSLProfile removes a client-ssl profile.
func (b *BigIP) DeleteClientSSLProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriClientSSL, name)
}

// ModifyClientSSLProfile allows you to change any attribute of a sever-ssl profile.
// Fields that can be modified are referenced in the VirtualClient struct.
func (b *BigIP) ModifyClientSSLProfile(name string, config *ClientSSLProfile) error {
	return b.put(config, uriLtm, uriProfile, uriClientSSL, name)
}
