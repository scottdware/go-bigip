package bigip

type SSLCertificateInfo struct {
	FileName    string `json:"file_name"`
	IsBundled   bool   `json:"is_bundled"`
	Certificate struct {
		CertInfo struct {
			ID    string      `json:"id"`
			Email interface{} `json:"email"`
		} `json:"cert_info"`
		ExpirationString string      `json:"expiration_string"`
		CertType         string      `json:"cert_type"`
		KeyType          string      `json:"key_type"`
		Version          int         `json:"version"`
		ExpirationDate   int         `json:"expiration_date"`
		SerialNumber     interface{} `json:"serial_number"`
		BitLength        int         `json:"bit_length"`
		Issuer           struct {
			DivisionName     string `json:"division_name"`
			StateName        string `json:"state_name"`
			LocalityName     string `json:"locality_name"`
			OrganizationName string `json:"organization_name"`
			CountryName      string `json:"country_name"`
			CommonName       string `json:"common_name"`
		} `json:"issuer"`
		Subject struct {
			DivisionName     string `json:"division_name"`
			StateName        string `json:"state_name"`
			LocalityName     string `json:"locality_name"`
			OrganizationName string `json:"organization_name"`
			CountryName      string `json:"country_name"`
			CommonName       string `json:"common_name"`
		} `json:"subject"`
	} `json:"certificate"`
}
