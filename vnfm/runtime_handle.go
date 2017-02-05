package vnfm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
	log "github.com/sirupsen/logrus"
)

type vnfmError struct {
	msg   string
	vnfr  *catalogue.VirtualNetworkFunctionRecord
	nsrID string
}

func (e *vnfmError) Error() string {
	return e.msg
}

type worker struct {
	*vnfm
	id int
}

func (wk *worker) allocateResources(
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	vimInstances map[string]*catalogue.VIMInstance,
	keyPairs []*catalogue.Key) (*catalogue.VirtualNetworkFunctionRecord, *vnfmError) {

	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle",
		"worker-id": wk.id,
		"vnfr-name": vnfr.Name,
	}).Debug("allocating resources for the VNFR")

	userData := wk.hnd.UserData()

	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle-allocate_resources",
		"worker-id": wk.id,
		"user-data": userData,
	}).Debug("will send to NFVO UserData")

	msg, err := messages.New(&messages.VNFMAllocateResources{
		VNFR:         vnfr,
		VIMInstances: vimInstances,
		Userdata:     userData,
		KeyPairs:     keyPairs,
	})
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle-allocate_resources",
			"worker-id": wk.id,
			"err":       err,
		}).Panicf("BUG")
	}

	nfvoResp, err := wk.cnl.NFVOExchange(msg)
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle-allocate_resources",
			"worker-id": wk.id,
			"err":       err,
		}).Error("exchange error")

		return nil, &vnfmError{
			msg:   "Unable to allocate Resources",
			nsrID: vnfr.ParentNsID,
			vnfr:  vnfr,
		}
	}

	if nfvoResp != nil {
		if nfvoResp.Action() == catalogue.ActionError {
			errorMessage := nfvoResp.Content().(*messages.OrError)

			wk.l.WithFields(log.Fields{
				"tag":          "worker-vnfm-handle-allocate_resources",
				"worker-id":    wk.id,
				"nfvo-err-msg": errorMessage.Message,
			}).Errorln("received error message from the NFVO")

			errVNFR := errorMessage.VNFR

			return nil, &vnfmError{
				msg:   fmt.Sprintf("Unable to allocate Resources. Reason: %s", errorMessage.Message),
				vnfr:  errVNFR,
				nsrID: vnfr.ParentNsID,
			}
		}

		message := nfvoResp.Content().(*messages.OrGeneric)
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle-allocate_resources",
			"worker-id": wk.id,
			"vnfr-name": message.VNFR.Name,
		}).Debug("received a VNFR from ALLOCATE")

		return message.VNFR, nil
	}

	return nil, &vnfmError{
		msg:   "received an empty message from the NFVO",
		nsrID: vnfr.ParentNsID,
		vnfr:  vnfr,
	}
}

func (wk *worker) handle(message messages.NFVMessage) error {
	content := message.Content()

	var reply messages.NFVMessage
	var err *vnfmError

	switch message.Action() {

	case catalogue.ActionConfigure:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = wk.handleConfigure(genericMessage)

	case catalogue.ActionError:
		errorMessage := content.(*messages.OrError)
		err = wk.handleError(errorMessage)

	case catalogue.ActionHeal:
		healMessage := content.(*messages.OrHealVNFRequest)
		reply, err = wk.handleHeal(healMessage)

	case catalogue.ActionInstantiate:
		instantiateMessage := content.(*messages.OrInstantiate)
		reply, err = wk.handleInstantiate(instantiateMessage)

	case catalogue.ActionInstantiateFinish:

	case catalogue.ActionModify:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = wk.handleModify(genericMessage)

	case catalogue.ActionScaleIn:
		scalingMessage := content.(*messages.OrScaling)
		err = wk.handleScaleIn(scalingMessage)

	case catalogue.ActionScaleOut:
		scalingMessage := content.(*messages.OrScaling)
		reply, err = wk.handleScaleOut(scalingMessage)

	// not implemented
	case catalogue.ActionScaling:

	case catalogue.ActionStart:
		startStopMessage := content.(*messages.OrStartStop)
		reply, err = wk.handleStart(startStopMessage)

	case catalogue.ActionStop:
		startStopMessage := content.(*messages.OrStartStop)
		reply, err = wk.handleStop(startStopMessage)

	// not implemented
	case catalogue.ActionReleaseResourcesFinish:

	case catalogue.ActionReleaseResources:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = wk.handleReleaseResources(genericMessage)

	case catalogue.ActionResume:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = wk.handleResume(genericMessage)

	case catalogue.ActionUpdate:
		updateMessage := content.(*messages.OrUpdate)
		reply, err = wk.handleUpdate(updateMessage)

	default:
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"action":    message.Action(),
		}).Warn("received unsupported action")

	}

	if err != nil {
		// send the error to the NFVO
		errorMsg, newMsgErr := messages.New(&messages.VNFMError{
			VNFR:  err.vnfr,
			NSRID: err.nsrID,
		})
		if newMsgErr != nil {
			wk.l.WithFields(log.Fields{
				"tag":       "worker-vnfm-handle",
				"worker-id": wk.id,
				"err":       err,
			}).Panic("BUG: shouldn't happen")
		}

		if sendErr := wk.cnl.NFVOSend(errorMsg); sendErr != nil {
			return fmt.Errorf("cannot send error message '%v' to the NFVO: %v", err, sendErr)
		}

		return err
	}

	if reply != nil {
		if reply.From() != messages.VNFM {
			wk.l.WithFields(log.Fields{
				"tag":           "worker-vnfm-handle",
				"worker-id":     wk.id,
				"msg-from-type": reply.From(),
			}).Panic("BUG: cannot send to the NFVO a message not intended to be received by it")
		}

		wk.l.WithFields(log.Fields{
			"tag":                "worker-vnfm-handle",
			"worker-id":          wk.id,
			"reply-action":       reply.Action(),
			"reply-content-type": reflect.TypeOf(reply.Content()).Name,
		}).Debug("sending reply to NFVO")

		if err := wk.cnl.NFVOSend(reply); err != nil {
			return fmt.Errorf("cannot send reply '%v' to the NFVO: %v", reply, err)
		}
	}

	return nil
}

func (wk *worker) handleConfigure(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	nfvMessage, err := messages.New(catalogue.ActionConfigure, &messages.VNFMGeneric{
		VNFR: genericMessage.VNFR,
	})
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleError(errorMessage *messages.OrError) *vnfmError {
	vnfr := errorMessage.VNFR
	nsrID := vnfr.ParentNsID

	wk.l.WithFields(log.Fields{
		"tag":            "worker-vnfm-handle",
		"worker-id":      wk.id,
		"nfvo-error-msg": errorMessage.Message,
	}).Errorf("received an error from the NFVO")

	if err := wk.hnd.HandleError(errorMessage.VNFR); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (wk *worker) handleHeal(healMessage *messages.OrHealVNFRequest) (messages.NFVMessage, *vnfmError) {
	vnfr := healMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := healMessage.VNFCInstance

	vnfrObtained, err := wk.hnd.Heal(vnfr, vnfcInstance, healMessage.Cause)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionHeal, &messages.VNFMHealed{
		VNFR:         vnfrObtained,
		VNFCInstance: vnfcInstance,
	})
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleInstantiate(instantiateMessage *messages.OrInstantiate) (messages.NFVMessage, *vnfmError) {
	extension := instantiateMessage.Extension

	wk.l.WithFields(log.Fields{
		"tag":        "worker-vnfm-handle",
		"worker-id":  wk.id,
		"extensions": extension,
	}).Debug("received extensions")

	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle",
		"worker-id": wk.id,
		"keys":      instantiateMessage.Keys,
	}).Debug("received keys")

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
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	resp, err := wk.cnl.NFVOExchange(msg)

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

	wk.l.WithFields(log.Fields{
		"tag":             "worker-vnfm-handle",
		"worker-id":       wk.id,
		"vnfr-hb_version": recvVNFR.HbVersion,
	}).Debug("received VNFR")

	if !wk.conf.Allocate {
		allocatedVNFR, err := wk.allocateResources(recvVNFR, vimInstanceChosen, instantiateMessage.Keys)
		if err != nil {
			return nil, err
		}

		recvVNFR = allocatedVNFR
	}

	var resultVNFR *catalogue.VirtualNetworkFunctionRecord

	if instantiateMessage.VNFPackage != nil {
		pkg := instantiateMessage.VNFPackage

		if pkg.ScriptsLink != "" {
			resultVNFR, err = wk.hnd.Instantiate(recvVNFR, pkg.ScriptsLink, vimInstances)
		} else {
			resultVNFR, err = wk.hnd.Instantiate(recvVNFR, pkg.Scripts, vimInstances)
		}
	} else {
		resultVNFR, err = wk.hnd.Instantiate(recvVNFR, nil, vimInstances)
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

func (wk *worker) handleModify(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := wk.hnd.Modify(vnfr, genericMessage.VNFRDependency)
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

func (wk *worker) handleReleaseResources(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := wk.hnd.Terminate(vnfr)
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

func (wk *worker) handleResume(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	vnfrDependency := genericMessage.VNFRDependency
	nsrID := vnfr.ParentNsID

	if actionForResume := wk.hnd.ActionForResume(vnfr, nil); actionForResume != catalogue.NoActionSpecified {
		resumedVNFR, err := wk.hnd.Resume(vnfr, nil, vnfrDependency)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		nfvMessage, err := messages.New(actionForResume, &messages.VNFMGeneric{
			VNFR: resumedVNFR,
		})
		if err != nil {
			wk.l.WithFields(log.Fields{
				"tag":       "worker-vnfm-handle",
				"worker-id": wk.id,
				"err":       err,
			}).Panic("BUG: shouldn't happen")
		}

		wk.l.WithFields(log.Fields{
			"tag":                    "worker-vnfm-handle",
			"worker-id":              wk.id,
			"vnfr-id":                vnfr.ID,
			"vnfr_dependency-target": vnfrDependency.Target,
			"resume-action":          actionForResume,
		}).Debug("Resuming VNFR")

		return nfvMessage, nil
	}

	return nil, nil
}

func (wk *worker) handleScaleIn(scalingMessage *messages.OrScaling) *vnfmError {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID

	vnfcInstanceToRemove := scalingMessage.VNFCInstance

	if _, err := wk.hnd.Scale(catalogue.ActionScaleIn, vnfr, vnfcInstanceToRemove, nil, nil); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (wk *worker) handleScaleOut(scalingMessage *messages.OrScaling) (messages.NFVMessage, *vnfmError) {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	component := scalingMessage.Component

	wk.l.WithFields(log.Fields{
		"tag":             "worker-vnfm-handle",
		"worker-id":       wk.id,
		"vnfr-hb_version": vnfr.HbVersion,
		"scaling_mode":    scalingMessage.Mode,
	}).Debug("received VNFR")

	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle",
		"worker-id": wk.id,
		"vnfc":      component,
	}).Info("Adding VNFComponent")

	var newVNFCInstance *catalogue.VNFCInstance
	if wk.conf.Allocate {
		newMsg, err := messages.New(&messages.VNFMScaling{
			VNFR:     vnfr,
			UserData: wk.hnd.UserData(),
		})

		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		var replyVNFR *catalogue.VirtualNetworkFunctionRecord

		switch content := newMsg.Content().(type) {
		case messages.OrGeneric:
			replyVNFR = content.VNFR
			wk.l.WithFields(log.Fields{
				"tag":                   "worker-vnfm-handle",
				"worker-id":             wk.id,
				"reply-vnfr-hb_version": replyVNFR.HbVersion,
			}).Debug("got reply VNFR")

		case messages.OrError:
			if err := wk.hnd.HandleError(content.VNFR); err != nil {
				return nil, &vnfmError{err.Error(), content.VNFR, nsrID}
			}

			return nil, nil

		default:
			replyVNFR = vnfr
		}

		if newVNFCInstance = replyVNFR.FindComponentInstance(component); newVNFCInstance == nil {
			return nil, &vnfmError{"no new VNFCInstance found. This should not happen.", replyVNFR, nsrID}
		}

		wk.l.WithFields(log.Fields{
			"tag":        "worker-vnfm-handle",
			"worker-id":  wk.id,
			"found-vnfc": newVNFCInstance.VNFComponent,
		}).Debug("VNFComponentInstance found")

		if strings.EqualFold(scalingMessage.Mode, "STANDBY") {
			newVNFCInstance.State = "STANDBY"
		}

		vnfr = replyVNFR
	} else {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
		}).Warn("wk.allocate is not set. No new VNFCInstance has been instantiated.")
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

	resultVNFR, err := wk.hnd.Scale(catalogue.ActionScaleOut, vnfr, newVNFCInstance, scripts, scalingMessage.Dependency)
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

func (wk *worker) handleStart(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error
	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = wk.hnd.Start(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = wk.hnd.StartVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStart, startStop)
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleStop(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error
	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = wk.hnd.Stop(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = wk.hnd.StopVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStop, startStop)
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleUpdate(updateMessage *messages.OrUpdate) (messages.NFVMessage, *vnfmError) {
	vnfr := updateMessage.VNFR
	nsrID := vnfr.ParentNsID
	script := updateMessage.Script

	replyVNFR, err := wk.hnd.UpdateSoftware(script, vnfr)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionUpdate, &messages.VNFMGeneric{
		VNFR: replyVNFR,
	})
	if err != nil {
		wk.l.WithFields(log.Fields{
			"tag":       "worker-vnfm-handle",
			"worker-id": wk.id,
			"err":       err,
		}).Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) work() {
	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle",
		"worker-id": wk.id,
	}).Debug("VNFM worker starting")

	// msgChan should be closed by the driver when exiting.
	for msg := range wk.msgChan {
		wk.l.WithFields(log.Fields{
			"tag":          "worker-vnfm-handle",
			"worker-id":    wk.id,
			"action":       msg.Action(),
			"content-type": reflect.TypeOf(msg.Content()).Name(),
		}).Debug("accepting new message")

		if err := wk.handle(msg); err != nil {
			wk.l.WithFields(log.Fields{
				"tag":       "worker-vnfm-handle",
				"worker-id": wk.id,
				"err":       err,
			}).Error("Handling error")
		} else {
			wk.l.WithFields(log.Fields{
				"tag":          "worker-vnfm-handle",
				"worker-id":    wk.id,
				"action":       msg.Action(),
				"content-type": reflect.TypeOf(msg.Content()).Name(),
			}).Debug("message successfully handled")
		}
	}

	wk.l.WithFields(log.Fields{
		"tag":       "worker-vnfm-handle",
		"worker-id": wk.id,
	}).Debug("VNFM worker exiting")

	wk.wg.Done()
}
