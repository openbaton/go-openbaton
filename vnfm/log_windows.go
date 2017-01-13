// +build windows

package vnfm

import (
	"io"
	"os"

	"github.com/shiena/ansicolor"
)

// Windows has issues with ANSI color codes, so stderr needs to be wrapped
// into an ANSI compatible Writer.
func terminalWriter() io.Writer {
	return ansicolor.NewAnsiColorWriter(os.Stderr)
}
