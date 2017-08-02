package bigip

// ServerSSLProfiles
// Documentation: https://devcentral.f5.com/wiki/iControlREST.APIRef_tm_ltm_profile_server-ssl.ashx

// ServerSSLProfiles contains a list of every server-ssl profile on the BIG-IP system.
type ServerSSLProfiles struct {
	ServerSSLProfiles []ServerSSLProfile `json:"items"`
}

// ServerSSLProfile contains information about each server-ssl profile. You can use all
// of these fields when modifying a server-ssl profile.
type ServerSSLProfile struct {
	Name                         string   `json:"name,omitempty"`
	Partition                    string   `json:"partition,omitempty"`
	FullPath                     string   `json:"fullPath,omitempty"`
	Generation                   int      `json:"generation,omitempty"`
	AlertTimeout                 string   `json:"alertTimeout,omitempty"`
	Authenticate                 string   `json:"authenticate,omitempty"`
	AuthenticateDepth            int      `json:"authenticateDepth,omitempty"`
	CaFile                       string   `json:"caFile,omitempty"`
	CacheSize                    int      `json:"cacheSize,omitempty"`
	CacheTimeout                 int      `json:"cacheTimeout,omitempty"`
	Cert                         string   `json:"cert,omitempty"`
	Chain                        string   `json:"chain,omitempty"`
	Ciphers                      string   `json:"ciphers,omitempty"`
	DefaultsFrom                 string   `json:"defaultsFrom,omitempty"`
	ExpireCertResponseControl    string   `json:"expireCertResponseControl,omitempty"`
	GenericAlert                 string   `json:"genericAlert,omitempty"`
	HandshakeTimeout             string   `json:"handshakeTimeout,omitempty"`
	Key                          string   `json:"key,omitempty"`
	ModSslMethods                string   `json:"modSslMethods,omitempty"`
	Mode                         string   `json:"mode,omitempty"`
	TmOptions                    []string `json:"tmOptions,omitempty"`
	Passphrase                   string   `json:"passphrase,omitempty"`
	PeerCertMode                 string   `json:"peerCertMode,omitempty"`
	ProxySsl                     string   `json:"proxySsl,omitempty"`
	RenegotiatePeriod            string   `json:"renegotiatePeriod,omitempty"`
	RenegotiateSize              string   `json:"renegotiateSize,omitempty"`
	Renegotiation                string   `json:"renegotiation,omitempty"`
	RetainCertificate            string   `json:"retainCertificate,omitempty"`
	SecureRenegotiation          string   `json:"secureRenegotiation,omitempty"`
	ServerName                   string   `json:"serverName,omitempty"`
	SessionMirroring             string   `json:"sessionMirroring,omitempty"`
	SessionTicket                string   `json:"sessionTicket,omitempty"`
	SniDefault                   string   `json:"sniDefault,omitempty"`
	SniRequire                   string   `json:"sniRequire,omitempty"`
	SslForwardProxy              string   `json:"sslForwardProxy,omitempty"`
	SslForwardProxyBypass        string   `json:"sslForwardProxyBypass,omitempty"`
	SslSignHash                  string   `json:"sslSignHash,omitempty"`
	StrictResume                 string   `json:"strictResume,omitempty"`
	UncleanShutdown              string   `json:"uncleanShutdown,omitempty"`
	UntrustedCertResponseControl string   `json:"untrustedCertResponseControl,omitempty"`
}

// ServerSSLProfiles returns a list of server-ssl profiles.
func (b *BigIP) ServerSSLProfiles() (*ServerSSLProfiles, error) {
	var serverSSLProfiles ServerSSLProfiles
	err, _ := b.getForEntity(&serverSSLProfiles, uriLtm, uriProfile, uriServerSSL)
	if err != nil {
		return nil, err
	}

	return &serverSSLProfiles, nil
}

// GetServerSSLProfile gets a server-ssl profile by name. Returns nil if the server-ssl profile does not exist
func (b *BigIP) GetServerSSLProfile(name string) (*ServerSSLProfile, error) {
	var serverSSLProfile ServerSSLProfile
	err, ok := b.getForEntity(&serverSSLProfile, uriLtm, uriProfile, uriServerSSL, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &serverSSLProfile, nil
}

// CreateServerSSLProfile creates a new server-ssl profile on the BIG-IP system.
func (b *BigIP) CreateServerSSLProfile(name string, parent string) error {
	config := &ServerSSLProfile{
		Name:         name,
		DefaultsFrom: parent,
	}

	return b.post(config, uriLtm, uriProfile, uriServerSSL)
}

// AddServerSSLProfile adds a new server-ssl profile on the BIG-IP system.
func (b *BigIP) AddServerSSLProfile(config *ServerSSLProfile) error {
	return b.post(config, uriLtm, uriProfile, uriServerSSL)
}

// DeleteServerSSLProfile removes a server-ssl profile.
func (b *BigIP) DeleteServerSSLProfile(name string) error {
	return b.delete(uriLtm, uriProfile, uriServerSSL, name)
}

// ModifyServerSSLProfile allows you to change any attribute of a sever-ssl profile.
// Fields that can be modified are referenced in the VirtualServer struct.
func (b *BigIP) ModifyServerSSLProfile(name string, config *ServerSSLProfile) error {
	return b.put(config, uriLtm, uriProfile, uriServerSSL, name)
}
