package plugin

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func (p *plug) closeLog() (err error) {
	if p.l.Out != os.Stderr {
		err = p.l.Out.(*os.File).Close()
		p.l.Out = os.Stderr
	}

	return
}

// initLogger opens and sets a log file
func (p *plug) openLog(defaultPath string) error {
	path := p.params.LogFile
	if path == "" {
		path = defaultPath
	}

	p.l = log.New()
	p.l.Level = p.params.LogLevel

	if path != "" {			
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("couldn't open the log file %s: %s", path, err.Error())
		}

		p.l.Out = file
	}

	return nil
}
