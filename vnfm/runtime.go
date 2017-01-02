package vnfm

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
)

var impls = make(map[string]channel.Driver)

func Register(name string, driver channel.Driver) {
	if _, ok := impls[name]; ok {
		panic(fmt.Sprintf("trying to register driver of type %T with already existing name '%s'", driver, name))
	}

	if driver == nil {
		panic("nil driver")
	}

	impls[name] = driver
}

type VNFM interface {
	Logger() *log.Logger
	Serve() error
	Stop() error
}

func New(implName string, handler Handler, config *config.Config) (VNFM, error) {
	if _, ok := impls[implName]; !ok {
		return nil, fmt.Errorf("no implementation available for %s. Have you forgot to import its package?", implName)
	}

	logger := log.New()

	logger.Formatter = &log.TextFormatter{
		ForceColors:   config.LogColors,
		DisableColors: !config.LogColors,
	}

	logger.Level = config.LogLevel

	if config.LogFile != "" {
		file, err := os.Open(config.LogFile)
		if err != nil {
			return nil, fmt.Errorf("couldn't open the log file %s: %s", config.LogFile, err.Error())
		}

		logger.Out = file
	}

	return &vnfm{
		hnd:      handler,
		implName: implName,
		conf:     config,
		l:        logger,
		quitChan: make(chan error, 1), // do not block on send
	}, nil
}

type vnfm struct {
	cnl      channel.Channel
	conf     *config.Config
	hnd      Handler
	implName string
	l        *log.Logger
	msgChan  <-chan messages.NFVMessage
	quitChan chan error
}

func (vnfm *vnfm) Logger() *log.Logger {
	return vnfm.l
}

func (vnfm *vnfm) Serve() (err error) {
	if vnfm.cnl, err = impls[vnfm.implName].Init(vnfm.conf, vnfm.l); err != nil {
		return
	}

	defer func() {
		r := recover()

		// If it's not stderr, it's the file we opened in New.
		if vnfm.l.Out != os.Stderr {
			vnfm.l.Out.(*os.File).Close()
		}

		// answering the channel signals Stop() that we're quitting
		vnfm.quitChan <- vnfm.cnl.Close()

		if r != nil {
			vnfm.l.WithFields(log.Fields{
				"tag":         "vnfm-serve-on_exit",
				"stack-trace": string(debug.Stack()),
			}).Panic(r)
		}
	}()

	if vnfm.msgChan, err = vnfm.cnl.NotifyReceived(); err != nil {
		return
	}

	vnfm.spawnWorkers()

	// wait for Stop()
	<-vnfm.quitChan

	return
}

func (vnfm *vnfm) SetLogger(log *log.Logger) {
	vnfm.l = log
}

func (vnfm *vnfm) Stop() error {
	select {
	case vnfm.quitChan <- nil:

	case <-time.After(time.Second):
		return errors.New("the VNFM is not listening")
	}

	select {
	case err := <-vnfm.quitChan:
		return err
	case <-time.After(1 * time.Minute):
		return errors.New("the VNFM refused to quit")
	}
}

func (vnfm *vnfm) spawnWorkers() {
	const NumWorkers = 5

	for i := 0; i < NumWorkers; i++ {
		go (&worker{vnfm, i}).spawn()
	}
}
