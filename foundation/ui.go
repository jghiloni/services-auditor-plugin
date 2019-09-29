package foundation

import (
	"os"

	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
)

// DefaultUI reads from stdin and writes to stdout
func DefaultUI() terminal.UI {
	return terminal.NewUI(
		os.Stdin,
		Writer,
		terminal.NewTeePrinter(Writer),
		trace.NewLogger(Writer, false, "false", ""),
	)
}
