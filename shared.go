package bigip

import (
	"fmt"
	"time"
)

const (
	uriClientSSL       = "client-ssl"
	uriDatagroup       = "data-group"
	uriHttp            = "http"
	uriHttpCompression = "http-compression"
	uriIRule           = "rule"
	uriInternal        = "internal"
	uriLtm             = "ltm"
	uriMonitor         = "monitor"
	uriNode            = "node"
	uriOneConnect      = "one-connect"
	uriPolicy          = "policy"
	uriPool            = "pool"
	uriPoolMember      = "members"
	uriProfile         = "profile"
	uriServerSSL       = "server-ssl"
	uriSnatPool        = "snatpool"
	uriTcp             = "tcp"
	uriUdp             = "udp"
	uriVirtual         = "virtual"
	uriVirtualAddress  = "virtual-address"
	uriGtm             = "gtm"
	uriWideIp          = "wideip"
	uriARecord         = "a"
	uriAAAARecord      = "aaaa"
	uriCNameRecord     = "cname"
	uriMXRecord        = "mx"
	uriNaptrRecord     = "naptr"
	uriSrvRecord       = "srv"
	uriPoolMembers     = "members"
	ENABLED            = "enable"
	DISABLED           = "disable"
	CONTEXT_SERVER     = "serverside"
	CONTEXT_CLIENT     = "clientside"
	CONTEXT_ALL        = "all"

	// Newer policy APIs have a draft-publish workflow that this library does not support.
	policyVersionSuffix = "?ver=11.5.1"
)

var cidr = map[string]string{
	"0":  "0.0.0.0",
	"1":  "128.0.0.0",
	"2":  "192.0.0.0",
	"3":  "224.0.0.0",
	"4":  "240.0.0.0",
	"5":  "248.0.0.0",
	"6":  "252.0.0.0",
	"7":  "254.0.0.0",
	"8":  "255.0.0.0",
	"9":  "255.128.0.0",
	"10": "255.192.0.0",
	"11": "255.224.0.0",
	"12": "255.240.0.0",
	"13": "255.248.0.0",
	"14": "255.252.0.0",
	"15": "255.254.0.0",
	"16": "255.255.0.0",
	"17": "255.255.128.0",
	"18": "255.255.192.0",
	"19": "255.255.224.0",
	"20": "255.255.240.0",
	"21": "255.255.248.0",
	"22": "255.255.252.0",
	"23": "255.255.254.0",
	"24": "255.255.255.0",
	"25": "255.255.255.128",
	"26": "255.255.255.192",
	"27": "255.255.255.224",
	"28": "255.255.255.240",
	"29": "255.255.255.248",
	"30": "255.255.255.252",
	"31": "255.255.255.254",
	"32": "255.255.255.255",
}

// GTMType handles the record types possible in the GTM as strings
type GTMType string

// GTM Record Types
const (
	ARecord     GTMType = uriARecord
	AAAARecord  GTMType = uriAAAARecord
	CNAMERecord GTMType = uriCNameRecord
	MXRecord    GTMType = uriMXRecord
	NAPTRRecord GTMType = uriNaptrRecord
	SRVRecord   GTMType = uriSrvRecord
)

const (
	uriShared       = "shared"
	uriLicensing    = "licensing"
	uriActivation   = "activation"
	uriRegistration = "registration"

	activationComplete   = "LICENSING_COMPLETE"
	activationInProgress = "LICENSING_ACTIVATION_IN_PROGRESS"
	activationFailed     = "LICENSING_FAILED"
	activationNeedEula   = "NEED_EULA_ACCEPT"
)

// https://devcentral.f5.com/wiki/iControl.Licensing_resource_API.ashx
type Activation struct {
	BaseRegKey            string   `json:"baseRegKey,omitempty"`
	AddOnKeys             []string `json:"addOnKeys,omitempty"`
	IsAutomaticActivation bool     `json:"isAutomaticActivation"`
	Status                string   `json:"status,omitempty"`
	LicenseText           *string  `json:"licenseText,omitempty"`
	ErrorText             *string  `json:"errorText,omitempty"`
	EulaText              *string  `json:"eulaText,omitempty"`
}

// https://devcentral.f5.com/wiki/iControl.Licensing_resource_API.ashx
type LicenseState struct {
	Vendor string `json:"vendor"`

	LicensedDateTime     string `json:"licensedDateTime"`
	LicensedVersion      string `json:"licensedVersion"`
	LicenseEndDateTime   string `json:"licenseEndDateTime"`
	LicenseStartDateTime string `json:"licenseStartDateTime"`

	RegistrationKey      string   `json:"registrationKey"`
	Dossier              string   `json:"dossier"`
	Authorization        string   `json:"authorization"`
	Usage                string   `json:"usage"`
	PlatformId           string   `json:"platformId"`
	AuthVers             string   `json:"authVers"`
	ServiceCheckDateTime string   `json:"serviceCheckDateTime"`
	MachineId            string   `json:"machineId"`
	ExclusivePlatform    []string `json:"exclusivePlatform"`

	ActiveModules   []string             `json:"activeModules"`
	OptionalModules []string             `json:"optionalModules"`
	FeatureFlags    []LicenseFeatureFlag `json:"featureFlags"`

	ExpiresInDays        string `json:"expiresInDays"`
	ExpiresInDaysMessage string `json:"expiresInDaysMessage"`
}

// Describes feature flags that are defined in licenses.
type LicenseFeatureFlag struct {
	FeatureName  string `json:"featureName"`
	FeatureValue string `json:"featureValue"`
}

// Gets the current activation status. Use after calling Activate. See the docs for more:
// https://devcentral.f5.com/wiki/iControl.Licensing_activation_APIs.ashx
func (b *BigIP) GetActivationStatus() (*Activation, error) {
	var a Activation
	err, _ := b.getForEntity(&a, uriShared, uriLicensing, uriActivation)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Sends the Activation to the activation endpoint. For documentation on how this works, see:
// https://devcentral.f5.com/wiki/iControl.Licensing_activation_APIs.ashx
func (b *BigIP) Activate(a Activation) error {
	return b.post(a, uriShared, uriLicensing, uriActivation)
}

// Returns the current license state.
func (b *BigIP) GetLicenseState() (*LicenseState, error) {
	var l LicenseState
	err, _ := b.getForEntity(&l, uriShared, uriLicensing, uriRegistration)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// Installs the given license.
func (b *BigIP) InstallLicense(licenseText string) error {
	r := map[string]string{"licenseText": licenseText}
	return b.put(r, uriShared, uriLicensing, uriRegistration)
}

// Automatically activate this registration key and install the resulting license.
// The BIG-IP must have access to the activation server for this to work.
func (b *BigIP) AutoLicense(regKey string, addOnKeys []string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	actreq := Activation{BaseRegKey: regKey, AddOnKeys: addOnKeys, IsAutomaticActivation: true}

	if err := b.Activate(actreq); err != nil {
		return err
	}

loop:
	for time.Now().Before(deadline) {
		actresp, err := b.GetActivationStatus()
		if err != nil {
			return err
		}

		if actresp.Status == activationInProgress {
			time.Sleep(1 * time.Second)
			continue
		}

		switch actresp.Status {
		case activationComplete:
			return b.InstallLicense(*actresp.LicenseText)
		case activationFailed:
			return fmt.Errorf("Licensing failed: %s", *actresp.ErrorText)
		case activationNeedEula:
			eula := *actresp.EulaText
			actreq.EulaText = &eula
			break loop
		default:
			return fmt.Errorf("Unknown licensing status: %s", actresp.Status)
		}
	}

	if actreq.EulaText == nil {
		return fmt.Errorf("Timed out after %s", timeout)
	}

	// Proceed with EULA acceptance
	if err := b.Activate(actreq); err != nil {
		return err
	}

	for time.Now().Before(deadline) {
		actresp, err := b.GetActivationStatus()
		if err != nil {
			return err
		}

		if actresp.Status == activationInProgress {
			time.Sleep(1 * time.Second)
			continue
		}

		switch actresp.Status {
		case activationComplete:
			return b.InstallLicense(*actresp.LicenseText)
		case activationNeedEula:
			return fmt.Errorf("Tried to accept EULA, but status is: %s", *actresp.ErrorText)
		case activationFailed:
			return fmt.Errorf("Licensing failed: %s", *actresp.ErrorText)
		}
		return fmt.Errorf("Unknown licensing status: %s", actresp.Status)
	}

	return fmt.Errorf("Timed out after %s", timeout)
}
