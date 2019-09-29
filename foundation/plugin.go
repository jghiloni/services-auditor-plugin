package foundation

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/cf/i18n"
	"code.cloudfoundry.org/cli/plugin"
	"github.com/blang/semver"
)

// Version is the plugin's version
var Version = "0.0.0"
var pluginVersion plugin.VersionType

func init() {
	v, _ := semver.ParseTolerant(Version)
	pluginVersion.Major = int(v.Major)
	pluginVersion.Minor = int(v.Minor)
	pluginVersion.Build = int(v.Patch)

	i18n.T = func(translationID string, args ...interface{}) string {
		if len(args) == 0 {
			return fmt.Sprintf("%s\n", translationID)
		}

		return fmt.Sprintf(translationID+"\n", args...)
	}
}

// GetMetadata returns usage and version info
func (p *ServiceAuditorPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:    "Service Auditor",
		Version: pluginVersion,
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "audit-services",
				HelpText: "Creates a report of all service instances in use",
				UsageDetails: plugin.Usage{
					Usage: "cf audit-services",
				},
			},
		},
	}
}

// Start binds the plugin
func (p *ServiceAuditorPlugin) Start() {
	plugin.Start(p)
}

// Run does the unit of work
func (p *ServiceAuditorPlugin) Run(cli plugin.CliConnection, args []string) {
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
		}
	}()

	// if uninstall is called, this will be some other value
	if args[0] != p.GetMetadata().Commands[0].Name {
		return
	}

	if !isUserAdmin(cli) {
		p.UI.Warn("You are not currently logged in as a user with administrator or read only administrator access. This may have unpredictable results.")
		if !p.UI.Confirm("Do you want to continue?") {
			p.UI.Warn("Bailing out")
			return
		}
	}

	p.UI.Say("Fetching information about all available services. This may take some time...\n")

	auditor := NewAuditor(cli)
	rows, err := auditor.Audit()
	if err != nil {
		p.UI.Failed(err.Error())
	}

	p.UI.Ok()
	table := p.UI.Table([]string{"Service Name", "Plan Name", "Instances", "Bindings", "Keys"})
	for i := range rows {
		table.Add(rows[i].TableRow()...)
	}
	err = table.Print()
	if err != nil {
		p.UI.Failed(err.Error())
	}
}

func isUserAdmin(cli plugin.CliConnection) bool {
	token, err := cli.AccessToken()
	if err != nil {
		return false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var accessToken AccessToken
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(&accessToken)
	if err != nil {
		return false
	}

	for _, scope := range accessToken.Scopes {
		if scope == "cloud_controller.admin" ||
			scope == "cloud_controller.admin_read_only" ||
			scope == "cloud_controller.global_auditor" {
			return true
		}
	}

	return false
}
