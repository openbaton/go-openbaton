// +build !windows

package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type logData struct{}

func (p *plug) deinitLogger() error {
	return p.l.Out.(*os.File).Close()
}

// initLogger opens and sets a log file on Unix platforms
func (p *plug) initLogger() error {
	pathArray := []string{"var", "log", p.params.Name + ".log"}
	path := filepath.Join(pathArray...)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return nil, fmt.Errorf("couldn't open the log file %s: %s", config.LogFile, err.Error())
	}

	p.l = log.New()
	p.l.Out = file
}
