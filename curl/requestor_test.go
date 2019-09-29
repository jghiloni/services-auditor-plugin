package curl_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/jghiloni/services-auditor-plugin/curl"
	"github.com/mitchellh/mapstructure"
)

func TestGetGood(t *testing.T) {
	testServer := getTestServer()
	defer testServer.Close()

	cli := getCliMock(testServer)
	cli.AccessTokenReturns(goodAuthToken, nil)

	requestor := curl.NewRequestor(cli)
	retMap, err := requestor.Get("/v2/info?hello=goodbye")
	if err != nil {
		t.Errorf("Expected error not to occur, got %q", err.Error())
		return
	}

	var v interface{}
	var ok bool
	if v, ok = retMap["key1"]; !ok {
		t.Error("Expected map to have key1, but did not")
		return
	}

	convertedV := make(map[string]interface{})
	err = mapstructure.Decode(v, &convertedV)
	if err != nil {
		t.Errorf("Expected error not to occur, got %q", err.Error())
		return
	}

	if _, ok = convertedV["key2"]; !ok {
		t.Error("Expected sub-map to have key2, but did not")
		return
	}
}

func TestGetUnauthorized(t *testing.T) {
	testServer := getTestServer()
	defer testServer.Close()

	cli := getCliMock(testServer)
	cli.AccessTokenReturns(badAuthToken, nil)

	requestor := curl.NewRequestor(cli)
	_, err := requestor.Get("/v2/info?hello=goodbye")
	if err == nil {
		t.Error("Expected error to occur, but none occured")
		return
	}

	if !strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
		t.Errorf("Expected to get an unauthorized error, got %q", err.Error())
		return
	}
}

func TestGetNotFound(t *testing.T) {
	testServer := getTestServer()
	defer testServer.Close()

	cli := getCliMock(testServer)
	cli.AccessTokenReturns(goodAuthToken, nil)

	requestor := curl.NewRequestor(cli)
	_, err := requestor.Get("/v1/info?hello=goodbye")
	if err == nil {
		t.Error("Expected error to occur, but none occured")
		return
	}

	if !strings.Contains(strings.ToLower(err.Error()), "not found") {
		t.Errorf("Expected to get a not found error, got %q", err.Error())
		return
	}
}

func getTestServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.EqualFold(authHeader, goodAuthToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(r.URL.Path, "/v2") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testJSON))
	}))
}

func getCliMock(server *httptest.Server) *pluginfakes.FakeCliConnection {
	cli := &pluginfakes.FakeCliConnection{}
	cli.ApiEndpointReturns(server.URL, nil)
	cli.IsSSLDisabledReturns(true, nil)

	return cli
}

const badAuthToken = `garbage-in-garbage-out`
const goodAuthToken = `bearer good-auth-token`
const testJSON = `{
	"key1": {
		"key2": [1,2,3]
	}
}`
