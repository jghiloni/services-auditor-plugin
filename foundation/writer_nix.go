// +build !windows

package foundation

import "os"

// Writer - a Linux-ish specific colorized output channel.
//    taken from https://github.com/cloudfoundry/cli/blob/master/cf/cmd/writer_unix.go
var Writer = os.Stdout
