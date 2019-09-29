package foundation_test

import (
	"testing"

	"github.com/jghiloni/services-auditor-plugin/foundation"
)

func TestServiceAuditor(t *testing.T) {
	cliConnection, server := getCLIConnection()
	defer server.Close()

	cliConnection.AccessTokenReturns("bearer abcd", nil)

	auditor := foundation.NewAuditor(cliConnection)
	rows, err := auditor.Audit()
	if err != nil {
		t.Errorf("Expected error not to occur, got %q", err.Error())
		return
	}

	if len(rows) != 101 {
		t.Errorf("Expected to get %d rows, got %d instead", 101, len(rows))
		return
	}

	for _, r := range rows {
		if r.ServiceName == "cleardb" && r.PlanName == "spark" {
			if r.Instances != 1 {
				t.Errorf("Expected cleardb/spark to have 1 instance, found %d", r.Instances)
				return
			}

			if r.Bindings != 3 {
				t.Errorf("Expected cleardb/spark to have 3 bindings, found %d", r.Bindings)
				return
			}

			if r.Keys != 1 {
				t.Errorf("Expected cleardb/spark to have 1 key, found %d", r.Keys)
				return
			}

			continue
		}

		if r.ServiceName == "cloudamqp" && r.PlanName == "lemur" {
			if r.Instances != 1 {
				t.Errorf("Expected cloudamqp/lemur to have 1 instance, found %d", r.Instances)
				return
			}

			if r.Bindings != 0 {
				t.Errorf("Expected cloudamqp/lemur to have 0 bindings, found %d", r.Bindings)
				return
			}

			if r.Keys != 1 {
				t.Errorf("Expected cloudamqp/lemur to have 1 key, found %d", r.Keys)
				return
			}

			continue
		}

		if r.ServiceName == "rediscloud" && r.PlanName == "30mb" {
			if r.Instances != 1 {
				t.Errorf("Expected rediscloud/30mb to have 1 instance, found %d", r.Instances)
				return
			}

			if r.Bindings != 1 {
				t.Errorf("Expected rediscloud/30mb to have 1 bindings, found %d", r.Bindings)
				return
			}

			if r.Keys != 0 {
				t.Errorf("Expected rediscloud/30mb to have 0 key, found %d", r.Keys)
				return
			}

			continue
		}

		if r.Instances != 0 {
			t.Errorf("Expected %s/%s to have 0 instance, found %d", r.ServiceName, r.PlanName, r.Instances)
			return
		}

		if r.Bindings != 0 {
			t.Errorf("Expected %s/%s to have 0 bindings, found %d", r.ServiceName, r.PlanName, r.Bindings)
			return
		}

		if r.Keys != 0 {
			t.Errorf("Expected %s/%s to have 0 key, found %d", r.ServiceName, r.PlanName, r.Keys)
			return
		}
	}
}
