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
	"log"
	"reflect"
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

func (r *RequestError) Error() error {
<<<<<<< HEAD
	if r.Message != "" {
=======
	if (r.Message != "") {
>>>>>>> origin/master
		return errors.New(r.Message)
	}
	return nil
}

// NewSession sets up our connection to the BIG-IP system.
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

	log.Println(options.Body)

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
		if res.Header["Content-Type"][0] == "application/json" {
			return data, b.checkError(data)
		}

		return data, errors.New(fmt.Sprintf("HTTP %d :: %s", res.StatusCode, string(data[:])))
	}

	return data, nil
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

// Perform a GET request and treat 404's as nil objects instead of errors.
func (b *BigIP) SafeGet(url string) ([]byte, error) {
	req := &APIRequest{
		Method:      "get",
		URL:         url,
		ContentType: "application/json",
	}

	resp, err := b.APICall(req)
	if err != nil {
		var reqError RequestError
		json.Unmarshal(resp, &reqError)
		if reqError.Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	return resp, nil
}
<<<<<<< HEAD
=======

// Helper to copy between transfer objects and model objects to hide the myriad of boolean representations
// in the iControlREST api. DTO fields can be tagged with bool:"yes|enabled|true" to set what true and false
// marshal to
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
		} else {
			return fmt.Errorf("Unknown type conversion %s -> %s", fromField.Kind(), toField.Kind())
		}
	}
	return nil
}

func toBoolString(b bool, trueStr, falseStr string) (string) {
	if b {
		return trueStr
	}
	return falseStr
}
>>>>>>> origin/master
