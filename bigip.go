// Package bigip interacts with F5 BIG-IP systems using the REST API.
package bigip

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	microToSeconds  = 1000000 // conversion factor
	maxTokenTimeout = 36000   // maximum token timeout in seconds
)

var defaultConfigOptions = &ConfigOptions{
	APICallTimeout: 60 * time.Second,
}

type ConfigOptions struct {
	APICallTimeout time.Duration
}

// BigIP is a container for our session state.
type BigIP struct {
	Host          string
	User          string
	Password      string
	Token         string    // if set, will be used instead of User/Password
	TokenExpiry   time.Time // the token expiration time
	Transport     *http.Transport
	ConfigOptions *ConfigOptions
	loginProvider string
	startTime     time.Time // token start time
}

// APIRequest builds our request before sending it to the server.
type APIRequest struct {
	Method      string
	URL         string
	Body        string
	ContentType string
}

// RequestError contains information about any error we get from a request.
type RequestError struct {
	Code       int      `json:"code,omitempty"`
	Message    string   `json:"message,omitempty"`
	ErrorStack []string `json:"errorStack,omitempty"`
}

// Upload contains information about a file upload status
type Upload struct {
	RemainingByteCount int64          `json:"remainingByteCount"`
	UsedChunks         map[string]int `json:"usedChunks"`
	TotalByteCount     int64          `json:"totalByteCount"`
	LocalFilePath      string         `json:"localFilePath"`
	TemporaryFilePath  string         `json:"temporaryFilePath"`
	Generation         int            `json:"generation"`
	LastUpdateMicros   int            `json:"lastUpdateMicros"`
}

// Error returns the error message.
func (r *RequestError) Error() error {
	if r.Message != "" {
		return errors.New(r.Message)
	}

	return nil
}

// NewSession sets up our connection to the BIG-IP system.
func NewSession(host, user, passwd string, configOptions *ConfigOptions) *BigIP {
	var url string
	if !strings.HasPrefix(host, "http") {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = host
	}
	if configOptions == nil {
		configOptions = defaultConfigOptions
	}
	return &BigIP{
		Host:     url,
		User:     user,
		Password: passwd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		ConfigOptions: configOptions,
	}
}

// NewTokenSession sets up our connection to the BIG-IP system, and
// instructs the session to use token authentication instead of Basic
// Auth. This is required when using an external authentication
// provider, such as Radius or Active Directory. loginProviderName is
// probably "tmos" but your environment may vary.
func NewTokenSession(host, user, passwd, loginProviderName string, configOptions *ConfigOptions) (b *BigIP, err error) {

	b = NewSession(host, user, passwd, configOptions)
	b.loginProvider = loginProviderName
	err = b.login()

	return
}

// APICall is used to query the BIG-IP web API.
func (b *BigIP) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{
		Transport: b.Transport,
		Timeout:   b.ConfigOptions.APICallTimeout,
	}
	var format string
	if strings.Contains(options.URL, "mgmt/") {
		format = "%s/%s"
	} else {
		format = "%s/mgmt/tm/%s"
	}
	url := fmt.Sprintf(format, b.Host, options.URL)
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)
	if b.Token != "" {
		req.Header.Set("X-F5-Auth-Token", b.Token)
	} else {
		req.SetBasicAuth(b.User, b.Password)
	}

	// fmt.Println("REQ -- ", options.Method, " ", url, " -- ", options.Body)

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		if res.Header.Get("Content-Type") == "application/json" {
			return data, b.checkError(data)
		}

		return data, errors.New(fmt.Sprintf("HTTP %d :: %s", res.StatusCode, string(data[:])))
	}

	// fmt.Println("Resp --", res.StatusCode, " -- ", string(data))
	return data, nil
}

// RefreshTokenSession refreshes the token expiration time by increasing
// token timeout by interal.
//
// If the token is already expired or if the above refresh fails, a new
// token is generated with a new login.
func (b *BigIP) RefreshTokenSession(interval time.Duration) error {
	if b.TokenExpiry.Sub(time.Now()) <= 0 {
		return b.login()
	}
	if err := b.increaseTokenTimout(interval); err != nil {
		fmt.Println(err)
		return b.login()
	}
	return nil
}

func (b *BigIP) iControlPath(parts []string) string {
	var buffer bytes.Buffer
	var lastPath int
	if strings.HasPrefix(parts[len(parts)-1], "?") {
		lastPath = len(parts) - 2
	} else {
		lastPath = len(parts) - 1
	}
	for i, p := range parts {
		buffer.WriteString(strings.Replace(p, "/", "~", -1))
		if i < lastPath {
			buffer.WriteString("/")
		}
	}
	return buffer.String()
}

//Generic delete
func (b *BigIP) delete(path ...string) error {
	req := &APIRequest{
		Method: "delete",
		URL:    b.iControlPath(path),
	}

	_, callErr := b.APICall(req)
	return callErr
}

func (b *BigIP) post(body interface{}, path ...string) error {
	return b.reqWithBody("post", body, path...)
}

func (b *BigIP) put(body interface{}, path ...string) error {
	return b.reqWithBody("put", body, path...)
}

func (b *BigIP) patch(body interface{}, path ...string) error {
	return b.reqWithBody("patch", body, path...)
}

func (b *BigIP) reqWithBody(method string, body interface{}, path ...string) error {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      method,
		URL:         b.iControlPath(path),
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	_, callErr := b.APICall(req)
	return callErr
}

//Get a url and populate an entity. If the entity does not exist (404) then the
//passed entity will be untouched and false will be returned as the second parameter.
//You can use this to distinguish between a missing entity or an actual error.
func (b *BigIP) getForEntity(e interface{}, path ...string) (error, bool) {
	req := &APIRequest{
		Method:      "get",
		URL:         b.iControlPath(path),
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		var reqError RequestError
		json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return nil, false
		}
		return err, false
	}

	err = json.Unmarshal(resp, e)
	if err != nil {
		return err, false
	}

	return nil, true
}

// checkError handles any errors we get from our API requests. It returns either the
// message of the error, if any, or nil.
func (b *BigIP) checkError(resp []byte) error {
	if len(resp) == 0 {
		return nil
	}

	var reqError RequestError

	err := json.Unmarshal(resp, &reqError)
	if err != nil {
		return errors.New(fmt.Sprintf("%s\n%s", err.Error(), string(resp[:])))
	}

	err = reqError.Error()
	if err != nil {
		return err
	}

	return nil
}

// jsonMarshal specifies an encoder with 'SetEscapeHTML' set to 'false' so that <, >, and & are not escaped. https://golang.org/pkg/encoding/json/#Marshal
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// Helper function to return a boolean pointer.
func Bool(b bool) *bool {
	return &b
}

// Helper to copy between transfer objects and model objects to hide the myriad of boolean representations
// in the iControlREST api. DTO fields can be tagged with bool:"yes|enabled|true" to set what true and false
// marshal to.
func marshal(to, from interface{}) error {
	toVal := reflect.ValueOf(to).Elem()
	fromVal := reflect.ValueOf(from).Elem()
	toType := toVal.Type()
	for i := 0; i < toVal.NumField(); i++ {
		toField := toVal.Field(i)
		toFieldType := toType.Field(i)
		fromField := fromVal.FieldByName(toFieldType.Name)
		if fromField.Interface() != nil && fromField.Kind() == toField.Kind() {
			toField.Set(fromField)
		} else if toField.Kind() == reflect.Bool && fromField.Kind() == reflect.String {
			switch fromField.Interface() {
			case "yes", "enabled", "true":
				toField.SetBool(true)
				break
			case "no", "disabled", "false", "":
				toField.SetBool(false)
				break
			default:
				return fmt.Errorf("Unknown boolean conversion for %s: %s", toFieldType.Name, fromField.Interface())
			}
		} else if _, ok := toField.Interface().(*bool); ok && fromField.Kind() == reflect.String {
			switch fromField.Interface() {
			case "yes", "enabled", "true":
				toField.Set(reflect.ValueOf(Bool(true)))
				break
			case "no", "disabled", "false":
				toField.Set(reflect.ValueOf(Bool(false)))
				break
			default:
				return fmt.Errorf("Unknown boolean conversion for %s: %s", toFieldType.Name, fromField.Interface())
			}
		} else if fromField.Kind() == reflect.Bool && toField.Kind() == reflect.String {
			tag := toFieldType.Tag.Get("bool")
			switch tag {
			case "yes":
				toField.SetString(toBoolString(fromField.Interface().(bool), "yes", "no"))
				break
			case "enabled":
				toField.SetString(toBoolString(fromField.Interface().(bool), "enabled", "disabled"))
				break
			case "true":
				toField.SetString(toBoolString(fromField.Interface().(bool), "true", "false"))
				break
			}
		} else if b, ok := fromField.Interface().(*bool); ok && toField.Kind() == reflect.String {
			if b == nil {
				continue
			}

			tag := toFieldType.Tag.Get("bool")
			switch tag {
			case "yes":
				toField.SetString(toBoolString(*b, "yes", "no"))
				break
			case "enabled":
				toField.SetString(toBoolString(*b, "enabled", "disabled"))
				break
			case "true":
				toField.SetString(toBoolString(*b, "true", "false"))
				break
			}
		} else {
			return fmt.Errorf("Unknown type conversion %s -> %s", fromField.Kind(), toField.Kind())
		}
	}
	return nil
}

func toBoolString(b bool, trueStr, falseStr string) string {
	if b {
		return trueStr
	}
	return falseStr
}

// Upload a file read from a Reader
func (b *BigIP) Upload(r io.Reader, size int64, path ...string) (*Upload, error) {
	client := &http.Client{
		Transport: b.Transport,
		Timeout:   b.ConfigOptions.APICallTimeout,
	}
	options := &APIRequest{
		Method:      "post",
		URL:         b.iControlPath(path),
		ContentType: "application/octet-stream",
	}
	var format string
	if strings.Contains(options.URL, "mgmt/") {
		format = "%s/%s"
	} else {
		format = "%s/mgmt/%s"
	}
	url := fmt.Sprintf(format, b.Host, options.URL)
	chunkSize := 512 * 1024
	var start, end int64
	for {
		// Read next chunk
		chunk := make([]byte, chunkSize)
		n, err := r.Read(chunk)
		if err != nil {
			return nil, err
		}
		end = start + int64(n)
		// Resize buffer size to number of bytes read
		if n < chunkSize {
			chunk = chunk[:n]
		}
		body := bytes.NewReader(chunk)
		req, _ := http.NewRequest(strings.ToUpper(options.Method), url, body)
		if b.Token != "" {
			req.Header.Set("X-F5-Auth-Token", b.Token)
		} else {
			req.SetBasicAuth(b.User, b.Password)
		}
		req.Header.Add("Content-Type", options.ContentType)
		req.Header.Add("Content-Range", fmt.Sprintf("%d-%d/%d", start, end-1, size))
		// Try to upload chunk
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		data, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode >= 400 {
			if res.Header.Get("Content-Type") == "application/json" {
				return nil, b.checkError(data)
			}

			return nil, fmt.Errorf("HTTP %d :: %s", res.StatusCode, string(data[:]))
		}
		defer res.Body.Close()
		var upload Upload
		err = json.Unmarshal(data, &upload)
		if err != nil {
			return nil, err
		}
		start = end
		if start >= size {
			// Final chunk was uploaded
			return &upload, err
		}
	}
}

// login requests a token.
func (b *BigIP) login() error {
	b.Token = ""
	b.startTime = time.Now()
	type authReq struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		LoginProviderName string `json:"loginProviderName"`
	}
	type authResp struct {
		Token struct {
			Token      string
			Expiration int `json:"expirationMicros"`
		}
	}

	auth := authReq{
		b.User,
		b.Password,
		b.loginProvider,
	}

	marshalJSON, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         "mgmt/shared/authn/login",
		Body:        string(marshalJSON),
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("unable to acquire authentication token")
	}

	var aresp authResp
	err = json.Unmarshal(resp, &aresp)
	if err != nil {
		return err
	}

	if aresp.Token.Token == "" {
		return fmt.Errorf("unable to acquire authentication token")
	}

	b.Token = aresp.Token.Token
	b.TokenExpiry = time.Unix(int64(aresp.Token.Expiration/microToSeconds), 0)

	return nil
}

// increaseTokenTimeout increases token timeout by interval.
//
// if it exceeds maxTokenTimeout an error is returned.
func (b *BigIP) increaseTokenTimout(interval time.Duration) error {
	if b.Token == "" {
		return errors.New("token refresh not possible - no token available")
	}
	newExpiry := time.Now().Add(interval)
	newTimeout := int(newExpiry.Sub(b.startTime)) / int(time.Second) // big ip token timeout is always relative to start time
	if newTimeout > maxTokenTimeout {
		return errors.New("maximum timeout exceeded")
	}

	type refreshReq struct {
		Timeout int `json:"timeout"`
	}
	type refreshResp struct {
		Token      string
		Expiration int `json:"expirationMicros"`
	}
	refreshJSON, err := json.Marshal(refreshReq{
		Timeout: newTimeout,
	})
	if err != nil {
		return err
	}

	req := &APIRequest{
		Method:      "patch",
		URL:         fmt.Sprintf("mgmt/shared/authz/tokens/%s", b.Token),
		Body:        string(refreshJSON),
		ContentType: "application/json",
	}
	resp, err := b.APICall(req)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("unable to refresh authentication token")
	}

	var rresp refreshResp
	err = json.Unmarshal(resp, &rresp)
	if err != nil {
		return err
	}
	if rresp.Expiration == 0 {
		return fmt.Errorf("unable to refresh authentication token")
	}
	b.TokenExpiry = time.Unix(int64(rresp.Expiration/microToSeconds), 0)
	return nil
}
