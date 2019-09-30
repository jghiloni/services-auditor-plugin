// +build windows

package foundation

import "github.com/fatih/color"

// Writer - a Windows specific colorized output channel.
//    taken from https://github.com/cloudfoundry/cli/blob/master/cf/cmd/writer_windows.go
var Writer = color.Output