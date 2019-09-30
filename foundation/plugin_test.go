package foundation_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/terminal/terminalfakes"
	"github.com/jghiloni/services-auditor-plugin/foundation"
)

func TestPluginAsAdmin(t *testing.T) {
	RegisterTestingT(t)

	ui := &terminalfakes.FakeUI{}
	ui.TableStub = func(headers []string) *terminal.UITable {
		return &terminal.UITable{
			UI:    ui,
			Table: terminal.NewTable(headers),
		}
	}

	plugin := foundation.NewAuditorPlugin(ui)

	cli, server := getCLIConnection()
	defer server.Close()

	cli.AccessTokenReturns(adminToken, nil)

	plugin.Run(cli, []string{plugin.GetMetadata().Commands[0].Name})

	Expect(ui.WarnCallCount()).To(Equal(0))
	Expect(ui.FailedCallCount()).To(Equal(0))
	Expect(ui.OkCallCount()).To(Equal(1))
	Expect(ui.SayCallCount()).To(Equal(2))

	_, args := ui.SayArgsForCall(1)
	Expect(args).To(HaveLen(1))

	argString, ok := args[0].(string)
	Expect(ok).To(BeTrue())

	tableLines := strings.Split(argString, "\n")
	Expect(tableLines).To(HaveLen(102))
}

func TestPluginAsAdminReadOnly(t *testing.T) {
	RegisterTestingT(t)

	ui := &terminalfakes.FakeUI{}
	plugin := foundation.NewAuditorPlugin(ui)

	cli, server := getCLIConnection()
	defer server.Close()

	cli.AccessTokenReturns(adminROToken, nil)

	plugin.Run(cli, []string{plugin.GetMetadata().Commands[0].Name})

	Expect(ui.WarnCallCount()).To(Equal(0))
}

func TestPluginAsAdminGlobalAuditor(t *testing.T) {
	RegisterTestingT(t)

	ui := &terminalfakes.FakeUI{}
	plugin := foundation.NewAuditorPlugin(ui)

	cli, server := getCLIConnection()
	defer server.Close()

	cli.AccessTokenReturns(adminGAToken, nil)

	plugin.Run(cli, []string{plugin.GetMetadata().Commands[0].Name})

	Expect(ui.WarnCallCount()).To(Equal(0))
}

func TestPluginAsNonAdminWithoutBail(t *testing.T) {
	RegisterTestingT(t)

	ui := &terminalfakes.FakeUI{}
	ui.ConfirmReturns(true)
	ui.TableStub = func(headers []string) *terminal.UITable {
		return &terminal.UITable{
			UI:    ui,
			Table: terminal.NewTable(headers),
		}
	}

	plugin := foundation.NewAuditorPlugin(ui)

	cli, server := getCLIConnection()
	defer server.Close()

	cli.AccessTokenReturns(nonAdminToken, nil)

	plugin.Run(cli, []string{plugin.GetMetadata().Commands[0].Name})
	Expect(ui.WarnCallCount()).To(Equal(1))
	Expect(ui.FailedCallCount()).To(Equal(0))
	Expect(ui.OkCallCount()).To(Equal(1))

	Expect(ui.SayCallCount()).To(Equal(2))

	_, args := ui.SayArgsForCall(1)
	Expect(args).To(HaveLen(1))

	argString, ok := args[0].(string)
	Expect(ok).To(BeTrue())

	tableLines := strings.Split(argString, "\n")
	Expect(tableLines).To(HaveLen(102))
}

func TestPluginAsNonAdminWithBail(t *testing.T) {
	RegisterTestingT(t)

	ui := &terminalfakes.FakeUI{}
	ui.ConfirmReturns(false)
	plugin := foundation.NewAuditorPlugin(ui)

	cli, server := getCLIConnection()
	defer server.Close()

	cli.AccessTokenReturns(nonAdminToken, nil)

	plugin.Run(cli, []string{plugin.GetMetadata().Commands[0].Name})
	Expect(ui.WarnCallCount()).To(Equal(2))
	Expect(ui.FailedCallCount()).To(Equal(0))
	Expect(ui.OkCallCount()).To(Equal(0))
	Expect(ui.SayCallCount()).To(Equal(0))
}

const adminToken = `foo.eyJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLmFkbWluIl19.bar`
const adminROToken = `foo.eyJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLmFkbWluX3JlYWRfb25seSJdfQ==.bar`
const adminGAToken = `foo.eyJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLmdsb2JhbF9hdWRpdG9yIl19.bar`

const nonAdminToken = `foo.eyJzY29wZSI6WyJzY2ltLnVzZXIiXX0=.bar`
