package foundation

import (
	"fmt"
	"net/url"

	"github.com/jghiloni/services-auditor-plugin/curl"
	"github.com/mitchellh/mapstructure"
)

// Audit does the thing
func (s *ServiceAuditor) Audit() (OutputRows, error) {
	services := make(map[string]interface{})
	err := s.buildResourceMap("/v2/services?q=active:true&results-per-page=100", services)
	if err != nil {
		return nil, err
	}

	servicePlans := make(map[string]interface{})
	for serviceGUID := range services {
		err = s.buildResourceMap(fmt.Sprintf("/v2/services/%s/service_plans?q=active:true&results-per-page=100", serviceGUID), servicePlans)
		if err != nil {
			return nil, err
		}
	}

	serviceInstances := make(map[string]interface{})
	siCountPerPlan := make(map[string]int)
	for servicePlanGUID := range servicePlans {
		siURL := fmt.Sprintf("/v2/service_plans/%s/service_instances", servicePlanGUID)
		siCountPerPlan[servicePlanGUID] = s.getCount(siURL)
		err = s.buildResourceMap(fmt.Sprintf("%s?results-per-page=100", siURL), serviceInstances)
		if err != nil {
			return nil, err
		}
	}

	rows := make([]OutputRow, 0, len(servicePlans))
	for serviceInstanceGUID, serviceInstance := range serviceInstances {
		var si ServiceInstance
		err = mapstructure.Decode(serviceInstance, &si)
		if err != nil {
			return nil, err
		}

		bindings := s.getCount(fmt.Sprintf("/v2/service_bindings?q=service_instance_guid:%s", serviceInstanceGUID))
		keys := s.getCount(fmt.Sprintf("/v2/service_keys?q=service_instance_guid:%s", serviceInstanceGUID))

		row := OutputRow{
			Bindings:  bindings,
			Keys:      keys,
			Instances: siCountPerPlan[si.ServicePlanGUID],
		}

		var sp ServicePlan
		servicePlan := servicePlans[si.ServicePlanGUID]
		err = mapstructure.Decode(servicePlan, &sp)
		if err != nil {
			return nil, err
		}
		row.PlanName = sp.Name

		service := services[sp.ServiceGUID]
		var svc Labeled
		err = mapstructure.Decode(service, &svc)
		if err != nil {
			return nil, err
		}
		row.ServiceName = svc.Label

		rows = append(rows, row)
	}

	for planGUID, siCount := range siCountPerPlan {
		if siCount == 0 {
			row := OutputRow{}

			var sp ServicePlan
			servicePlan := servicePlans[planGUID]
			err = mapstructure.Decode(servicePlan, &sp)
			if err != nil {
				return nil, err
			}
			row.PlanName = sp.Name

			service := services[sp.ServiceGUID]
			var svc Labeled
			err = mapstructure.Decode(service, &svc)
			if err != nil {
				return nil, err
			}
			row.ServiceName = svc.Label

			rows = append(rows, row)
		}
	}

	return rows, nil
}

func (s *ServiceAuditor) buildResourceMap(url string, resources map[string]interface{}) error {
	for {
		servicePage, err := s.getPage(url)
		if err != nil {
			return err
		}

		for _, resource := range servicePage.Resources {
			if err != nil {
				return err
			}

			var id Identifiable
			err = mapstructure.Decode(resource.Metadata, &id)
			if err != nil {
				return err
			}

			resources[id.GUID] = resource.Entity
		}

		url = servicePage.NextURL
		if url == "" {
			break
		}
	}

	return nil
}

func (s *ServiceAuditor) getPage(url string) (Page, error) {
	r := curl.NewRequestor(s.cli)

	m, e := r.Get(url)
	if e != nil {
		return Page{}, e
	}

	var p Page
	e = mapstructure.Decode(m, &p)
	if e != nil {
		return Page{}, e
	}

	return p, nil
}

func (s *ServiceAuditor) getCount(urlString string) int {
	u, e := url.ParseRequestURI(urlString)
	if e != nil {
		return 0
	}

	q := u.Query()
	q.Add("results-per-page", "1")

	p, e := s.getPage(fmt.Sprintf("%s?%s", u.EscapedPath(), q.Encode()))
	if e != nil {
		return 0
	}

	return p.TotalResults
}
