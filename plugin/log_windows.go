// +build windows

package plugin

import (
	"golang.org/x/sys/windows/svc/eventlog"

	log "github.com/sirupsen/logrus"
)

type logData *eventLogHook

// eventLogHook allows logrus to log to Windows EventLog
type eventLogHook struct {
	elog *eventlog.Log
	src  string
}

func (p *plug) deinitLogger() error {
	if err := (*eventLogHook)(p.e).Close(); err != nil {
		return err
	}

	return p.closeLogFile()
}

// initLogger creates a logger with an EventLog hook (requires admin privileges)
func (p *plug) initLogger() error {
	if err := eventlog.InstallAsEventCreate(p.params.Name,
		eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {

		return err
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

	p.l = log.New()
	p.l.Hooks.Add(lh)

	return p.openLogFile("") // no default
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
