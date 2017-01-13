// +build !windows

package vnfm

import (
    "io"
    "os"
)

// UNIX terminals and friends are assumed to be fine with ANSI color codes.
func terminalWriter() io.Writer {
    return os.Stderr
}