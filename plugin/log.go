/*
 *  Copyright (c) 2017 Open Baton (http://openbaton.org)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package plugin

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func (p *plug) closeLog() (err error) {
	if file, ok := p.l.Out.(*os.File); ok {
		err = file.Close()
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

	if path != "" && path != "-" {
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if err != nil {
			return fmt.Errorf("couldn't open the log file %s: %s", path, err.Error())
		}

		p.l.Out = file
	}

	return nil
}
