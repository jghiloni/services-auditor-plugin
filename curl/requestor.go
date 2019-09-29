package curl

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// Requestor will make calls to the CAPI
type Requestor struct {
	apiEndpoint string
	cli         plugin.CliConnection
	hc          *http.Client
}

// NewRequestor will set up a requestor object based on the CliConnection
func NewRequestor(cli plugin.CliConnection) *Requestor {
	api, _ := cli.ApiEndpoint()
	skipSSL, _ := cli.IsSSLDisabled()

	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipSSL,
			},
		},
	}

	return &Requestor{
		apiEndpoint: api,
		hc:          hc,
		cli:         cli,
	}
}

// Get does an API Get with authorization
func (r *Requestor) Get(uri string) (map[string]interface{}, error) {
	baseURL, err := url.Parse(r.apiEndpoint)
	if err != nil {
		return nil, err
	}

	requestURI, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	baseURL.Path = requestURI.Path
	baseURL.RawQuery = requestURI.RawQuery

	request, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	token, err := r.cli.AccessToken()
	if err != nil {
		return nil, err
	}

	token = strings.TrimSpace(token)
	if !strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = fmt.Sprintf("Bearer %s", token)
	}

	request.Header.Add("Authorization", token)
	response, err := r.hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("Error fetching resource: %s", response.Status)
	}

	responseMap := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&responseMap)

	return responseMap, err
}
