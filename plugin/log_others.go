// +build !windows

package plugin

import (
	"path/filepath"
)

type logData struct{}

func (p *plug) deinitLogger() error {
	return p.closeLog()
}

// initLogger opens and sets a log file on Unix platforms
func (p *plug) initLogger() error {
	pathArray := []string{"/", "var", "log", p.params.Name + ".log"}
	defaultPath := filepath.Join(pathArray...)

	return p.openLog(defaultPath)
}
