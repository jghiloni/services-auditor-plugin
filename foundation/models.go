package foundation

import (
	"strconv"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/plugin"
)

// Metadata represents info about a CF resource
type Metadata struct {
	GUID      string `mapstructure:"guid"`
	URL       string `mapstructure:"url"`
	CreatedAt string `mapstructure:"created_at"`
	UpdatedAt string `mapstructure:"updated_at"`
}

// Resource is a CF resource
type Resource struct {
	Metadata Metadata               `mapstructure:"metadata"`
	Entity   map[string]interface{} `mapstructure:"entity"`
}

// Page is a collection of CF resources
type Page struct {
	TotalResults int        `mapstructure:"total_results"`
	TotalPages   int        `mapstructure:"total_pages"`
	PreviousURL  string     `mapstructure:"prev_url"`
	NextURL      string     `mapstructure:"next_url"`
	Resources    []Resource `mapstructure:"resources"`
}

// Identifiable has a GUID field
type Identifiable struct {
	GUID string `mapstructure:"guid"`
}

// Named has a Name field
type Named struct {
	Name string `mapstructure:"name"`
}

// Labeled has a Label field
type Labeled struct {
	Label string `mapstructure:"label"`
}

// ServiceInstance has information about the service plan from which it was created
type ServiceInstance struct {
	ServicePlanGUID string `mapstructure:"service_plan_guid"`
}

// ServicePlan has information about the service that presents it
type ServicePlan struct {
	Named       `mapstructure:",squash"`
	ServiceGUID string `mapstructure:"service_guid"`
}

// AccessToken has information about a user's scope
type AccessToken struct {
	Scopes []string `json:"scope"`
}

// OutputRow provides the structure of the info to be sent to the terminal
type OutputRow struct {
	ServiceName string
	PlanName    string
	Instances   int
	Keys        int
	Bindings    int
}

// TableRow will convert the row into something consumable by the CF CLI's UI Table
func (o OutputRow) TableRow() []string {
	return []string{
		o.ServiceName,
		o.PlanName,
		strconv.Itoa(o.Instances),
		strconv.Itoa(o.Bindings),
		strconv.Itoa(o.Keys),
	}
}

// OutputRows wraps []OutputRow so we can support sort.Interface
type OutputRows []OutputRow

// Len returns the number of items in the slice
func (o OutputRows) Len() int { return len(o) }

// Swap swaps the items
func (o OutputRows) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o OutputRows) Less(i, j int) bool {
	if o[i].ServiceName != o[j].ServiceName {
		return o[i].ServiceName < o[j].ServiceName
	}

	if o[i].PlanName != o[j].PlanName {
		return o[i].PlanName < o[j].PlanName
	}

	return true
}

// ServiceAuditorPlugin will search a given foundation for service instances, bindings, and keys
type ServiceAuditorPlugin struct {
	UI terminal.UI
}

// ServiceAuditor will do the thing
type ServiceAuditor struct {
	cli plugin.CliConnection
}

// NewAuditorPlugin returns an auditor with the given UI
func NewAuditorPlugin(ui terminal.UI) *ServiceAuditorPlugin {
	return &ServiceAuditorPlugin{
		UI: ui,
	}
}

// NewAuditor creates a new auditor
func NewAuditor(cli plugin.CliConnection) *ServiceAuditor {
	return &ServiceAuditor{
		cli: cli,
	}
}
