package vnfm

import (
	"errors"
	"fmt"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/mcilloni/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

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

	if config.LogFile != "" {
		file, err := os.Open(config.LogFile)
		if err != nil {
			return nil, fmt.Errorf("couldn't open the log file %s: %s", config.LogFile, err.Error())
		}

		logger.Out = file
	}

	return &vnfm{
		hnd:      handler,
		conf:     config,
		l:        logger,
		quitChan: make(chan struct{}),
	}, nil
}

var impls map[string]channel.Driver

type vnfm struct {
	cnl      channel.Channel
	conf     *config.Config
	hnd      Handler
	implName string
	l        *log.Logger
	msgChan  <-chan messages.NFVMessage
	quitChan chan struct{}
}

func (vnfm *vnfm) Logger() *log.Logger {
	return vnfm.l
}

func (vnfm *vnfm) Serve() error {
	var err error
	if vnfm.cnl, err = impls[vnfm.implName].Init(vnfm.conf, vnfm.l); err != nil {
		return err
	}

	defer func() {
		r := recover()

		// If it's not stderr, it's the file we opened in New.
		if vnfm.l.Out != os.Stderr {
			vnfm.l.Out.(*os.File).Close()
		}

		if err := vnfm.cnl.Close(); err != nil {
			vnfm.l.Errorln(err)
		}

		if r != nil {
			vnfm.l.Panicln(r)
		}
	}()

	if vnfm.msgChan, err = vnfm.cnl.NotifyReceived(); err != nil {
		return err
	}

MainLoop:
	for {
		select {
		case msg := <-vnfm.msgChan:
			if err := vnfm.handle(msg); err != nil {
				vnfm.l.Errorln(err)
			}

		case <-vnfm.quitChan:
			break MainLoop

		default:

		}
	}

	return nil
}

func (vnfm *vnfm) SetLogger(log *log.Logger) {
	vnfm.l = log
}

func (vnfm *vnfm) Stop() error {
	vnfm.quitChan <- struct{}{}

	select {
	case <-vnfm.quitChan:
		return nil
	case <-time.After(time.Second):
		return errors.New("the vnfm refused to quit")
	}
}

type vnfmError struct {
	msg   string
	vnfr  *catalogue.VirtualNetworkFunctionRecord
	nsrID string
}

func (e *vnfmError) Error() string {
	return e.msg
}

func (vnfm *vnfm) allocateResources(
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	vimInstances map[string]*catalogue.VIMInstance,
	keyPairs []*catalogue.Key) (*catalogue.VirtualNetworkFunctionRecord, *vnfmError) {

	userData := vnfm.hnd.UserData()
	vnfm.l.Debugf("Userdata sent to NFVO: %s\n", userData)

	msg, err := messages.New(&messages.VNFMAllocateResources{
		VNFR:         vnfr,
		VIMInstances: vimInstances,
		Userdata:     userData,
		KeyPairs:     keyPairs,
	})
	if err != nil {
		vnfm.l.Panicf("BUG: %v\n", err)
	}

	nfvoResp, err := vnfm.cnl.NFVOExchange(msg)
	if err != nil {
		vnfm.l.Errorln(err.Error())
		return nil, &vnfmError{
			msg:   "Not able to allocate Resources",
			nsrID: vnfr.ParentNsID,
			vnfr:  vnfr,
		}
	}

	if nfvoResp != nil {
		if nfvoResp.Action() == catalogue.ActionError {
			errorMessage := nfvoResp.Content().(*messages.OrError)

			vnfm.l.Errorln(errorMessage.Message)

			errVNFR := errorMessage.VNFR

			return nil, &vnfmError{
				msg:   fmt.Sprintf("Not able to allocate Resources because: %s\n", errorMessage.Message),
				vnfr:  errVNFR,
				nsrID: vnfr.ParentNsID,
			}
		}

		message := nfvoResp.Content().(*messages.OrGeneric)
		vnfm.l.Debugf("Received from ALLOCATE: %s\n", message.VNFR)

		return message.VNFR, nil
	}

	return nil, &vnfmError{
		msg:   "received an empty message from NFVO",
		nsrID: vnfr.ParentNsID,
		vnfr:  vnfr,
	}
}

func (vnfm *vnfm) handle(message messages.NFVMessage) error {

	vnfm.l.Debugf("vnfm: Received Message: '%s'\n", message.Action())

	content := message.Content()

	var reply messages.NFVMessage
	var err *vnfmError

	switch message.Action() {
	case catalogue.ActionScaleIn:
		scalingMessage := content.(*messages.OrScaling)
		err = vnfm.handleScaleIn(scalingMessage)

	case catalogue.ActionScaleOut:
		scalingMessage := content.(*messages.OrScaling)
		reply, err = vnfm.handleScaleOut(scalingMessage)

	// not implemented
	case catalogue.ActionScaling:

	case catalogue.ActionError:
		errorMessage := content.(*messages.OrError)
		err = vnfm.handleError(errorMessage)

	case catalogue.ActionModify:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = vnfm.handleModify(genericMessage)

	case catalogue.ActionReleaseResources:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = vnfm.handleReleaseResources(genericMessage)

	case catalogue.ActionInstantiate:
		instantiateMessage := content.(*messages.OrInstantiate)
		reply, err = vnfm.handleInstantiate(instantiateMessage)

	// not implemented
	case catalogue.ActionReleaseResourcesFinish:

	case catalogue.ActionUpdate:
		updateMessage := content.(*messages.OrUpdate)
		reply, err = vnfm.handleUpdate(updateMessage)

	case catalogue.ActionHeal:
		healMessage := content.(*messages.OrHealVNFRequest)
		reply, err = vnfm.handleHeal(healMessage)

	case catalogue.ActionInstantiateFinish:

	case catalogue.ActionConfigure:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = vnfm.handleConfigure(genericMessage)

	case catalogue.ActionStart:
		startStopMessage := content.(*messages.OrStartStop)
		reply, err = vnfm.handleStart(startStopMessage)

	case catalogue.ActionStop:
		startStopMessage := content.(*messages.OrStartStop)
		reply, err = vnfm.handleStop(startStopMessage)

	case catalogue.ActionResume:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = vnfm.handleResume(genericMessage)

	default:
		vnfm.l.Warnf("received unsupported action '%s'\n", message.Action())

	}

	if err != nil {
		vnfm.l.Errorln(err.Error())

		errorMsg, err := messages.New(&messages.VNFMError{
			VNFR:  err.vnfr,
			NSRID: err.nsrID,
		})
		if err != nil {
			vnfm.l.Panicf("BUG: shouldn't happen: %v\n", err)
		}

		if err := vnfm.cnl.NFVOSend(errorMsg); err != nil {
			vnfm.l.Errorf("cannot send error message to the NFVO: %v\n", err)
		}
	} else {
		if reply != nil {
			if reply.From() != messages.VNFM {
				vnfm.l.Panicln("BUG: cannot send to the NFVO a message not intended to be received by it")
			}
			vnfm.l.Debugf("sending action: '%s' and a content '%T' to NFVO", reply.Action(), reply.Content())

			if err := vnfm.cnl.NFVOSend(reply); err != nil {
				vnfm.l.Errorf("cannot send a reply to the NFVO: %v\n", err)
			}
		}
	}

	return err
}

func (vnfm *vnfm) handleConfigure(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	nfvMessage, err := messages.New(catalogue.ActionConfigure, &messages.VNFMGeneric{
		VNFR: genericMessage.VNFR,
	})
	if err != nil {
		vnfm.l.Panicf("BUG: shouldn't happen: %v\n", err)
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleError(errorMessage *messages.OrError) *vnfmError {
	vnfr := errorMessage.VNFR
	nsrID := vnfr.ParentNsID

	vnfm.l.Errorf("received an error from the NFVO: %s\n", errorMessage.Message)

	if err := vnfm.hnd.HandleError(errorMessage.VNFR); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (vnfm *vnfm) handleHeal(healMessage *messages.OrHealVNFRequest) (messages.NFVMessage, *vnfmError) {
	vnfr := healMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := healMessage.VNFCInstance

	vnfrObtained, err := vnfm.hnd.Heal(vnfr, vnfcInstance, healMessage.Cause)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionHeal, &messages.VNFMHealed{
		VNFR:         vnfrObtained,
		VNFCInstance: vnfcInstance,
	})
	if err != nil {
		vnfm.l.Panicf("BUG: shouln't happen: %v\n", err)
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleInstantiate(instantiateMessage *messages.OrInstantiate) (messages.NFVMessage, *vnfmError) {
	extension := instantiateMessage.Extension

	vnfm.l.Debugf("received extensions: %v\n", extension)
	vnfm.l.Debugf("received keys: %v\n", instantiateMessage.Keys)

	vimInstances := instantiateMessage.VIMInstances

	vnfr, err := catalogue.NewVNFR(
		instantiateMessage.VNFD,
		instantiateMessage.VNFDFlavour.FlavourKey,
		instantiateMessage.VLRs,
		instantiateMessage.Extension,
		vimInstances)

	msg, err := messages.New(catalogue.ActionGrantOperation, &messages.VNFMGeneric{
		VNFR: vnfr,
	})
	if err != nil {
		vnfm.l.Panicf("BUG: should not happen: %v\n", err)
	}

	resp, err := vnfm.cnl.NFVOExchange(msg)

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, vnfr.ParentNsID}
	}

	respContent, ok := resp.Content().(*messages.OrGrantLifecycleOperation)
	if !ok {
		return nil, &vnfmError{
			msg:   fmt.Sprintf("expected OrGrantLifecycleOperation, got %T", resp.Content()),
			nsrID: vnfr.ParentNsID,
			vnfr:  vnfr,
		}
	}

	recvVNFR := respContent.VNFR
	vimInstanceChosen := respContent.VDUVIM

	vnfm.l.Debugf("VERSION IS: %d\n", recvVNFR.HbVersion)

	if vnfm.conf.Allocate {
		allocatedVNFR, err := vnfm.allocateResources(recvVNFR, vimInstanceChosen, instantiateMessage.Keys)
		if err != nil {
			return nil, err
		}

		recvVNFR = allocatedVNFR
	}

	var resultVNFR *catalogue.VirtualNetworkFunctionRecord

	if instantiateMessage.VNFPackage != nil {
		pkg := instantiateMessage.VNFPackage

		if pkg.ScriptsLink != "" {
			resultVNFR, err = vnfm.hnd.Instantiate(recvVNFR, pkg.ScriptsLink, vimInstances)
		} else {
			resultVNFR, err = vnfm.hnd.Instantiate(recvVNFR, pkg.Scripts, vimInstances)
		}
	} else {
		resultVNFR, err = vnfm.hnd.Instantiate(recvVNFR, nil, vimInstances)
	}

	if err != nil {
		return nil, &vnfmError{
			msg:   err.Error(),
			nsrID: recvVNFR.ParentNsID,
			vnfr:  recvVNFR,
		}
	}

	nfvMessage, err := messages.New(catalogue.ActionInstantiate, &messages.VNFMInstantiate{
		VNFR: resultVNFR,
	})
	if err != nil {
		return nil, &vnfmError{err.Error(), resultVNFR, resultVNFR.ParentNsID}
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleModify(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := vnfm.hnd.Modify(vnfr, genericMessage.VNFRDependency)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionModify, &messages.VNFMGeneric{
		VNFR: resultVNFR,
	})
	if err != nil {
		return nil, &vnfmError{err.Error(), resultVNFR, nsrID}
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleReleaseResources(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := vnfm.hnd.Terminate(vnfr)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionReleaseResources, &messages.VNFMGeneric{
		VNFR: resultVNFR,
	})
	if err != nil {
		return nil, &vnfmError{err.Error(), resultVNFR, nsrID}
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleResume(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	vnfrDependency := genericMessage.VNFRDependency
	nsrID := vnfr.ParentNsID

	if actionForResume := vnfm.hnd.ActionForResume(vnfr, nil); actionForResume != catalogue.NoActionSpecified {
		resumedVNFR, err := vnfm.hnd.Resume(vnfr, nil, vnfrDependency)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		nfvMessage, err := messages.New(actionForResume, &messages.VNFMGeneric{
			VNFR: resumedVNFR,
		})
		if err != nil {
			vnfm.l.Panicf("BUG: shouln't happen: %v\n", err)
		}

		vnfm.l.Debugf("Resuming vnfr '%d' with dependency target: '%s' for action: '%s'\n",
			vnfr.ID, vnfrDependency.Target, actionForResume)

		return nfvMessage, nil
	}

	return nil, nil
}

func (vnfm *vnfm) handleScaleIn(scalingMessage *messages.OrScaling) *vnfmError {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID

	vnfcInstanceToRemove := scalingMessage.VNFCInstance

	if _, err := vnfm.hnd.Scale(catalogue.ActionScaleIn, vnfr, vnfcInstanceToRemove, nil, nil); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (vnfm *vnfm) handleScaleOut(scalingMessage *messages.OrScaling) (messages.NFVMessage, *vnfmError) {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	component := scalingMessage.Component

	vnfm.l.Debugf("HB_VERSION == %d\n", vnfr.HbVersion)
	vnfm.l.Infof("Adding VNFComponent: %v\n", component)
	vnfm.l.Debugf("The mode is: %s\n", scalingMessage.Mode)

	var newVNFCInstance *catalogue.VNFCInstance
	if vnfm.conf.Allocate {
		newMsg, err := messages.New(&messages.VNFMScaling{
			VNFR:     vnfr,
			UserData: vnfm.hnd.UserData(),
		})

		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		var replyVNFR *catalogue.VirtualNetworkFunctionRecord

		switch content := newMsg.Content().(type) {
		case messages.OrGeneric:
			replyVNFR = content.VNFR
			vnfm.l.Debugf("HB_VERSION == %d\n", replyVNFR.HbVersion)

		case messages.OrError:
			if err := vnfm.hnd.HandleError(content.VNFR); err != nil {
				return nil, &vnfmError{err.Error(), content.VNFR, nsrID}
			}

			return nil, nil

		default:
			replyVNFR = vnfr
		}

		if newVNFCInstance = replyVNFR.FindComponentInstance(component); newVNFCInstance == nil {
			return nil, &vnfmError{"no new VNFCInstance found. This should not happen.", replyVNFR, nsrID}
		}

		vnfm.l.Debugf("VNFComponentInstance FOUND : %v\n", newVNFCInstance.VNFComponent)

		if strings.EqualFold(scalingMessage.Mode, "STANDBY") {
			newVNFCInstance.State = "STANDBY"
		}

		vnfr = replyVNFR
	} else {
		vnfm.l.Warnln("vnfm.allocate is not set. No new VNFCInstance has been instantiated.")
	}

	var scripts interface{}
	switch {
	case scalingMessage.VNFPackage == nil:
		scripts = []*catalogue.Script{}

	case scalingMessage.VNFPackage.ScriptsLink != "":
		scripts = scalingMessage.VNFPackage.ScriptsLink

	default:
		scripts = scalingMessage.VNFPackage.Scripts
	}

	resultVNFR, err := vnfm.hnd.Scale(catalogue.ActionScaleOut, vnfr, newVNFCInstance, scripts, scalingMessage.Dependency)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionScaled, &messages.VNFMScaled{
		VNFR:         resultVNFR,
		VNFCInstance: newVNFCInstance,
	})

	if err != nil {
		return nil, &vnfmError{err.Error(), resultVNFR, nsrID}
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleStart(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error
	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = vnfm.hnd.Start(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = vnfm.hnd.StartVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStart, startStop)
	if err != nil {
		vnfm.l.Panicf("BUG: shouln't happen: %v\n", err)
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleStop(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error
	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = vnfm.hnd.Stop(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = vnfm.hnd.StopVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStop, startStop)
	if err != nil {
		vnfm.l.Panicf("BUG: shouln't happen: %v\n", err)
	}

	return nfvMessage, nil
}

func (vnfm *vnfm) handleUpdate(updateMessage *messages.OrUpdate) (messages.NFVMessage, *vnfmError) {
	vnfr := updateMessage.VNFR
	nsrID := vnfr.ParentNsID
	script := updateMessage.Script

	replyVNFR, err := vnfm.hnd.UpdateSoftware(script, vnfr)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionUpdate, &messages.VNFMGeneric{
		VNFR: replyVNFR,
	})
	if err != nil {
		vnfm.l.Panicf("BUG: shouldn't happen: %v\n", err)
	}

	return nfvMessage, nil
}
