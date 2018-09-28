package bigip

import (
	"fmt"
	"time"
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
