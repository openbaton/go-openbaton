package vnfmsdk

import (
	"fmt"
	"errors"
	"strings"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
)

type worker struct {
	l        *logging.Logger
	handler  HandlerVnfm
	Allocate bool
	Channel  *amqp.Channel
}

type vnfmError struct {
	msg   string
	vnfr  *catalogue.VirtualNetworkFunctionRecord
	nsrID string
}

func (wk *worker) handleConfigure(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	nfvMessage, err := messages.New(catalogue.ActionConfigure, &messages.VNFMGeneric{
		VNFR: genericMessage.VNFR,
	})
	if err != nil {
		wk.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleError(errorMessage *messages.OrError) *vnfmError {
	vnfr := errorMessage.VNFR
	nsrID := vnfr.ParentNsID

	wk.l.Errorf("received an error from the NFVO")

	if err := wk.handler.HandleError(errorMessage.VNFR); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (wk *worker) handleHeal(healMessage *messages.OrHealVNFRequest) (messages.NFVMessage, *vnfmError) {
	vnfr := healMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := healMessage.VNFCInstance

	vnfrObtained, err := wk.handler.Heal(vnfr, vnfcInstance, healMessage.Cause)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionHeal, &messages.VNFMHealed{
		VNFR:         vnfrObtained,
		VNFCInstance: vnfcInstance,
	})
	if err != nil {
		wk.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleInstantiate(instantiateMessage *messages.OrInstantiate) (messages.NFVMessage, *vnfmError) {

	wk.l.Debug("received extensions: %v", instantiateMessage.Extension)

	wk.l.Debug("received keys: %v", instantiateMessage.Keys)

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
		wk.l.Panic("BUG: shouldn't happen")
	}

	resp, err := wk.executeRpc("vnfm.nfvo.actions.reply", msg)

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

	wk.l.Debug("received VNFR")

	if !wk.Allocate {
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
			resultVNFR, err = wk.handler.Instantiate(recvVNFR, pkg.ScriptsLink, vimInstances)
		} else {
			resultVNFR, err = wk.handler.Instantiate(recvVNFR, pkg.Scripts, vimInstances)
		}
	} else {
		resultVNFR, err = wk.handler.Instantiate(recvVNFR, nil, vimInstances)
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
func (worker *worker) executeRpc(queue string, message messages.NFVMessage) (messages.NFVMessage, error) {
	msgs, corrId, err := sdk.ExecuteRpc(queue, message, worker.Channel, worker.l)
	if err != nil {
		return nil, err
	}
	for d := range msgs {
		if corrId == d.CorrelationId {
			return messages.Unmarshal(d.Body, messages.NFVO)
		}
	}

	return nil, errors.New("no answer!")
}

func (wk *worker) handleModify(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := wk.handler.Modify(vnfr, genericMessage.VNFRDependency)
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

	resultVNFR, err := wk.handler.Terminate(vnfr)
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

	if actionForResume := wk.handler.ActionForResume(vnfr, nil); actionForResume != catalogue.NoActionSpecified {
		resumedVNFR, err := wk.handler.Resume(vnfr, nil, vnfrDependency)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		nfvMessage, err := messages.New(actionForResume, &messages.VNFMGeneric{
			VNFR: resumedVNFR,
		})
		if err != nil {
			wk.l.Panic("BUG: shouldn't happen")
		}

		wk.l.Debug("Resuming VNFR")

		return nfvMessage, nil
	}

	return nil, nil
}

func (wk *worker) handleScaleIn(scalingMessage *messages.OrScaling) *vnfmError {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID

	vnfcInstanceToRemove := scalingMessage.VNFCInstance

	if _, err := wk.handler.Scale(catalogue.ActionScaleIn, vnfr, vnfcInstanceToRemove, nil, nil); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (wk *worker) handleScaleOut(scalingMessage *messages.OrScaling) (messages.NFVMessage, *vnfmError) {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	component := scalingMessage.Component

	wk.l.Debug("received VNFR")

	wk.l.Info("Adding VNFComponent")

	var newVNFCInstance *catalogue.VNFCInstance
	if !wk.Allocate {
		newMsg, err := messages.New(&messages.VNFMScaling{
			VNFR:     vnfr,
			UserData: wk.handler.UserData(),
		})

		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		respMsg, err := wk.executeRpc("vnfm.nfvo.actions.reply", newMsg)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		var replyVNFR *catalogue.VirtualNetworkFunctionRecord

		switch content := respMsg.Content().(type) {
		case *messages.OrGeneric:
			replyVNFR = content.VNFR
			wk.l.Debug("got reply VNFR")

		case *messages.OrError:
			if err := wk.handler.HandleError(content.VNFR); err != nil {
				return nil, &vnfmError{err.Error(), content.VNFR, nsrID}
			}

			return nil, nil

		default:
			wk.l.Warning("got a weird message on reply to SCALING")

			replyVNFR = vnfr
		}

		if newVNFCInstance = replyVNFR.FindComponentInstance(component); newVNFCInstance == nil {
			return nil, &vnfmError{"no new VNFCInstance found. This should not happen.", replyVNFR, nsrID}
		}

		wk.l.Debug("VNFComponentInstance found")

		if strings.EqualFold(scalingMessage.Mode, "STANDBY") {
			newVNFCInstance.State = "STANDBY"
		}

		vnfr = replyVNFR
	} else {
		wk.l.Warning("wk.allocate is set. No new VNFCInstance has been instantiated by the NFVO.")
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

	resultVNFR, err := wk.handler.Scale(catalogue.ActionScaleOut, vnfr, newVNFCInstance, scripts, scalingMessage.Dependency)
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
		startStop.VNFR, err = wk.handler.Start(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = wk.handler.StartVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStart, startStop)
	if err != nil {
		wk.l.Panic("BUG: shouldn't happen")
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
		startStop.VNFR, err = wk.handler.Stop(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = wk.handler.StopVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStop, startStop)
	if err != nil {
		wk.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) handleUpdate(updateMessage *messages.OrUpdate) (messages.NFVMessage, *vnfmError) {
	vnfr := updateMessage.VNFR
	nsrID := vnfr.ParentNsID
	script := updateMessage.Script

	replyVNFR, err := wk.handler.UpdateSoftware(script, vnfr)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionUpdate, &messages.VNFMGeneric{
		VNFR: replyVNFR,
	})
	if err != nil {
		wk.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (wk *worker) allocateResources(
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	vimInstances map[string]*catalogue.BaseVimInstance,
	keyPairs []*catalogue.Key) (*catalogue.VirtualNetworkFunctionRecord, *vnfmError) {

	wk.l.Debug("allocating resources for the VNFR")

	userData := wk.handler.UserData()

	wk.l.Debug("will send to NFVO UserData")

	msg, err := messages.New(&messages.VNFMAllocateResources{
		VNFR:         vnfr,
		VIMInstances: vimInstances,
		Userdata:     userData,
		KeyPairs:     keyPairs,
	})
	if err != nil {
		wk.l.Panicf("BUG")
	}

	nfvoResp, err := wk.executeRpc("vnfm.nfvo.actions.reply", msg)
	if err != nil {
		wk.l.Error("exchange error")

		return nil, &vnfmError{
			msg:   "Unable to allocate Resources",
			nsrID: vnfr.ParentNsID,
			vnfr:  vnfr,
		}
	}

	if nfvoResp != nil {
		if nfvoResp.Action() == catalogue.ActionError {
			errorMessage := nfvoResp.Content().(*messages.OrError)

			wk.l.Error("received error message from the NFVO")

			errVNFR := errorMessage.VNFR

			return nil, &vnfmError{
				msg:   fmt.Sprintf("Unable to allocate Resources. Reason: %s", errorMessage.Message),
				vnfr:  errVNFR,
				nsrID: vnfr.ParentNsID,
			}
		}

		message := nfvoResp.Content().(*messages.OrGeneric)
		wk.l.Debug("received a VNFR from ALLOCATE")

		return message.VNFR, nil
	}

	return nil, &vnfmError{
		msg:   "received an empty message from the NFVO",
		nsrID: vnfr.ParentNsID,
		vnfr:  vnfr,
	}
}
