// +build !windows

package plugin

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type logData struct{}

func (p *plug) deinitLogger() error {
	return p.closeLogFile()
}

// initLogger opens and sets a log file on Unix platforms
func (p *plug) initLogger() error {
	pathArray := []string{"var", "log", p.params.Name + ".log"}
	defaultPath := filepath.Join(pathArray...)

	p.l = log.New()
	p.l.Level = p.params.LogLevel

	return p.openLogFile(defaultPath)
}
