// +build !windows

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
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type logData struct{}

func (p *plug) deinitLogger() error {
	return p.closeLog()
}

// initLogger opens and sets a log file on Unix platforms
func (p *plug) initLogger() error {
	pathArray := []string{"/", "var", "log", "openbaton", p.params.Type + "-plugin.log"}
	defaultPath := filepath.Join(pathArray...)

	if err := p.openLog(defaultPath); err != nil {
		return err
	}

	if p.params.LogFile == "-" {
		p.l.Formatter = &log.TextFormatter{
			DisableTimestamp: !p.params.Timestamps,
			FullTimestamp: p.params.Timestamps,
		}
	}

	return nil
}
