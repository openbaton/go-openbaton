package plugin

import (
	"fmt"
	"os"
)

func (p *plug) closeLogFile() error {
	if p.l.Out != os.Stderr {
		return p.l.Out.(*os.File).Close()
	}

	return nil
}

// initLogger opens and sets a log file
func (p *plug) openLogFile(defaultPath string) error {
	path := p.params.LogFile
	if path == "" {
		if defaultPath == "" {
			return nil
		}

		path = defaultPath
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("couldn't open the log file %s: %s", path, err.Error())
	}

	p.l.Out = file

	return nil
}