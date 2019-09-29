package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/jghiloni/services-auditor-plugin/foundation"
)

func main() {
	plugin.Start(foundation.NewAuditorPlugin(foundation.DefaultUI()))
}
