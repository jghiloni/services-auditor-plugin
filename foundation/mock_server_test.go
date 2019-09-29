package foundation_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
)

func getCLIConnection() (*pluginfakes.FakeCliConnection, *httptest.Server) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		replacer := strings.NewReplacer("/", "_", "?", "_", ":", "_", "=", "_", "&", "_")
		uri, err := url.PathUnescape(r.URL.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fileName := filepath.Join("testdata", replacer.Replace(uri))

		fp, err := os.Open(fileName)
		if err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusBadGateway)
			return
		}
		defer fp.Close()

		body, err := ioutil.ReadAll(fp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(body)
	}))

	cliConnection := &pluginfakes.FakeCliConnection{}
	cliConnection.ApiEndpointReturns(server.URL, nil)
	cliConnection.IsSSLDisabledReturns(true, nil)

	return cliConnection, server
}
