package plugin

import (
	log "github.com/sirupsen/logrus"
)

type Params struct {
	BrokerIP, LogFile, Name, Password, Username string
	Workers, Port                               int
}

type Plugin interface {
	Logger() *log.Logger
	Serve() error
	Stop() error
	Type() string
}

func NewDriver(driver Driver, p *Params) (Plugin, error) {
	return nil, nil
}

type plug struct {
	l      *log.Logger
	e      logData
	params *Params
	rh     reqHandler
}

func (p *plug) Logger() *log.Logger {
	return p.l
}

func (p *plug) Serve() error {
	return nil
}

func (p *plug) Stop() error {
	return nil
}

func (p *plug) Type() string {
	return p.rh.Type()
}

type reqHandler interface {
	Handle(call string, args interface{}) (interface{}, error)
	Type() string
}
