// +build windows

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
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
)

type logData *eventLogHook

// eventLogHook allows logrus to log to Windows EventLog
type eventLogHook struct {
	elog *eventlog.Log
	src  string
}

func (p *plug) deinitLogger() error {
	if p.params.LogFile == "" {
		if err := (*eventLogHook)(p.e).Close(); err != nil {
			return err
		}
	}

	return p.closeLog()
}

// initLogger creates a logger with an EventLog hook (requires admin privileges)
func (p *plug) initLogger() error {
	// no default
	if err := p.openLog(""); err != nil {
		return err
	}

	if p.params.LogFile == "-" {
		p.l.Formatter = &log.TextFormatter{
			DisableColors:    false,
			DisableTimestamp: !p.params.Timestamps,
			ForceColors:      true,
			FullTimestamp:    p.params.Timestamps,
			TimestampFormat:  time.RFC3339Nano,
		}

		p.l.Out = ansicolor.NewAnsiColorWriter(os.Stdout)
	}

	// do not enable the event logger if the logfile is present.

	if p.params.LogFile != "" {
		return nil
	}

	// try to install the event; if it fails because it already exists, try to
	// remove it and install again
	for {
		err := eventlog.InstallAsEventCreate(p.params.Name,
			eventlog.Error|eventlog.Warning|eventlog.Info)

		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "registry key already exists") {
			return err
		}

		if err := eventlog.Remove(p.params.Name); err != nil {
			return err
		}
	}

	el, err := eventlog.Open(p.params.Name)
	if err != nil {
		return err
	}

	lh := &eventLogHook{
		elog: el,
		src:  p.params.Name,
	}

	p.e = logData(lh)

	p.l.Hooks.Add(lh)

	return nil
}

// Close closes the logger and uninstalls the source
func (h *eventLogHook) Close() error {
	if err := h.elog.Close(); err != nil {
		return err
	}

	h.elog = nil

	return eventlog.Remove(h.src)
}

// Fire logs an entry to the EventLog.
func (h *eventLogHook) Fire(entry *log.Entry) error {
	if h.elog == nil {
		return nil
	}

	message, err := entry.String()
	if err != nil {
		return err
	}

	switch entry.Level {
	case log.PanicLevel:
		fallthrough
	case log.FatalLevel:
		fallthrough
	case log.ErrorLevel:
		return h.elog.Error(1, message)

	case log.WarnLevel:
		return h.elog.Warning(10, message)

	case log.InfoLevel:
		fallthrough
	case log.DebugLevel:
		return h.elog.Info(100, message)

	default:
		panic("unsupported level in hooks")
	}
}

// Levels returns the supported logging levels.
func (eventLogHook) Levels() []log.Level {
	return log.AllLevels
}
