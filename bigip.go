// Package bigip interacts with F5 BIG-IP systems using the REST API.
package bigip

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// BigIP is a container for our session state.
type BigIP struct {
	Host      string
	User      string
	Password  string
	Transport *http.Transport
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

// NewServer sets up our connection to the BIG-IP system.
func NewSession(host, user, passwd string) *BigIP {
	return &BigIP{
		Host:     host,
		User:     user,
		Password: passwd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// APICall is used to query the BIG-IP web API.
func (b *BigIP) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{Transport: b.Transport}
	url := fmt.Sprintf("https://%s/mgmt/tm/%s", b.Host, options.URL)
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)
	req.SetBasicAuth(b.User, b.Password)

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	return data, nil
}

// checkError handles any errors we get from our API requests. It returns either the
// message of the error, if any, or nil.
func (b *BigIP) checkError(resp []byte) error {
	var reqError RequestError
	err := json.Unmarshal(resp, &reqError)
	if err != nil {
		return err
	}

	if reqError.Message != "" {
		return errors.New(reqError.Message)
	}

	return nil
}
