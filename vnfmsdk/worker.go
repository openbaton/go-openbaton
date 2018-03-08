package vnfmsdk

import (
	"fmt"
	"strings"

	"github.com/op/go-logging"
	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/streadway/amqp"
)

//The worker struct allows the VNFM SDK to invoke implementation specific of VNFMs
type worker struct {
	l          *logging.Logger
	handler    HandlerVnfm
	Allocate   bool
	Connection *amqp.Connection
}

type vnfmError struct {
	msg   string
	vnfr  *catalogue.VirtualNetworkFunctionRecord
	nsrID string
}

func (worker *worker) handleConfigure(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	nfvMessage, err := messages.New(catalogue.ActionConfigure, &messages.VNFMGeneric{
		VNFR: genericMessage.VNFR,
	})
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (worker *worker) handleError(errorMessage *messages.OrError) *vnfmError {
	vnfr := errorMessage.VNFR
	nsrID := vnfr.ParentNsID

	worker.l.Errorf("received an error from the NFVO")

	if err := worker.handler.HandleError(errorMessage.VNFR); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (worker *worker) handleHeal(healMessage *messages.OrHealVNFRequest) (messages.NFVMessage, *vnfmError) {
	vnfr := healMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := healMessage.VNFCInstance

	vnfrObtained, err := worker.handler.Heal(vnfr, vnfcInstance, healMessage.Cause)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionHeal, &messages.VNFMHealed{
		VNFR:         vnfrObtained,
		VNFCInstance: vnfcInstance,
	})
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (worker *worker) handleInstantiate(instantiateMessage *messages.OrInstantiate) (messages.NFVMessage, *vnfmError) {

	worker.l.Debug("received extensions: %v", instantiateMessage.Extension)

	worker.l.Debug("received keys: %v", instantiateMessage.Keys)

	vimInstances := instantiateMessage.VIMInstances

	var flavorKey string

	if instantiateMessage.VNFDFlavour != nil {
		flavorKey = instantiateMessage.VNFDFlavour.FlavourKey
	} else {
		flavorKey = ""
	}

	vnfr, err := catalogue.NewVNFR(
		instantiateMessage.VNFD,
		flavorKey,
		instantiateMessage.VLRs,
		instantiateMessage.Extension,
		vimInstances)

	msg, err := messages.New(catalogue.ActionGrantOperation, &messages.VNFMGeneric{
		VNFR: vnfr,
	})
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	resp, err := worker.executeRpc("vnfm.nfvo.actions.reply", msg)

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

	worker.l.Debug("Received VNFR after GRANT_OPERATION")

	if !worker.Allocate {
		allocatedVNFR, err := worker.allocateResources(recvVNFR, vimInstanceChosen, instantiateMessage.Keys)
		if err != nil {
			return nil, err
		}

		recvVNFR = allocatedVNFR
	}

	var resultVNFR *catalogue.VirtualNetworkFunctionRecord

	if instantiateMessage.VNFPackage != nil {
		pkg := instantiateMessage.VNFPackage

		if pkg.ScriptsLink != "" {
			resultVNFR, err = worker.handler.Instantiate(recvVNFR, pkg.ScriptsLink, vimInstances)
		} else {
			resultVNFR, err = worker.handler.Instantiate(recvVNFR, pkg.Scripts, vimInstances)
		}
	} else {
		resultVNFR, err = worker.handler.Instantiate(recvVNFR, nil, vimInstances)
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
	body, err := sdk.Rpc(
		queue,
		message,
		worker.Connection,
		worker.l)

	if err != nil {
		return nil, err
	}
	return messages.Unmarshal(body, messages.NFVO)
}

func (worker *worker) handleModify(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := worker.handler.Modify(vnfr, genericMessage.VNFRDependency)
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

func (worker *worker) handleReleaseResources(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := worker.handler.Terminate(vnfr)
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

func (worker *worker) handleResume(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfmError) {
	vnfr := genericMessage.VNFR
	vnfrDependency := genericMessage.VNFRDependency
	nsrID := vnfr.ParentNsID

	if actionForResume := worker.handler.ActionForResume(vnfr, nil); actionForResume != catalogue.NoActionSpecified {
		resumedVNFR, err := worker.handler.Resume(vnfr, nil, vnfrDependency)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		nfvMessage, err := messages.New(actionForResume, &messages.VNFMGeneric{
			VNFR: resumedVNFR,
		})
		if err != nil {
			worker.l.Panic("BUG: shouldn't happen")
		}

		worker.l.Debug("Resuming VNFR")

		return nfvMessage, nil
	}

	return nil, nil
}

func (worker *worker) handleScaleIn(scalingMessage *messages.OrScaling) *vnfmError {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID

	vnfcInstanceToRemove := scalingMessage.VNFCInstance

	if _, _, err := worker.handler.Scale(scalingMessage.VIMInstance, catalogue.ActionScaleIn, vnfr, vnfcInstanceToRemove, nil, nil); err != nil {
		return &vnfmError{err.Error(), vnfr, nsrID}
	}

	return nil
}

func (worker *worker) handleScaleOut(scalingMessage *messages.OrScaling) (messages.NFVMessage, *vnfmError) {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	component := scalingMessage.Component

	worker.l.Debug("received VNFR")

	worker.l.Info("Adding VNFComponent")

	var newVNFCInstance *catalogue.VNFCInstance
	if !worker.Allocate {
		newMsg, err := messages.New(&messages.VNFMScaling{
			VNFR:     vnfr,
			UserData: worker.handler.UserData(),
		})

		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		respMsg, err := worker.executeRpc("vnfm.nfvo.actions.reply", newMsg)
		if err != nil {
			return nil, &vnfmError{err.Error(), vnfr, nsrID}
		}

		var replyVNFR *catalogue.VirtualNetworkFunctionRecord

		switch content := respMsg.Content().(type) {
		case *messages.OrGeneric:
			replyVNFR = content.VNFR
			worker.l.Debug("got reply VNFR")

		case *messages.OrError:
			if err := worker.handler.HandleError(content.VNFR); err != nil {
				return nil, &vnfmError{err.Error(), content.VNFR, nsrID}
			}

			return nil, nil

		default:
			worker.l.Warning("got a weird message on reply to SCALING")

			replyVNFR = vnfr
		}

		if newVNFCInstance = replyVNFR.FindComponentInstance(component); newVNFCInstance == nil {
			return nil, &vnfmError{"no new VNFCInstance found. This should not happen.", replyVNFR, nsrID}
		}

		worker.l.Debug("VNFComponentInstance found")

		if strings.EqualFold(scalingMessage.Mode, "STANDBY") {
			newVNFCInstance.State = "STANDBY"
		}

		vnfr = replyVNFR
	} else {
		worker.l.Warning("wk.allocate is set. No new VNFCInstance has been instantiated by the NFVO.")
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

	var err error
	var resultVNFR *catalogue.VirtualNetworkFunctionRecord

	if !worker.Allocate {
		resultVNFR, newVNFCInstance, err = worker.handler.Scale(scalingMessage.VIMInstance, catalogue.ActionScaleOut, vnfr, newVNFCInstance, scripts, scalingMessage.Dependency)
	} else {
		resultVNFR, newVNFCInstance, err = worker.handler.Scale(scalingMessage.VIMInstance, catalogue.ActionScaleOut, vnfr, component, scripts, scalingMessage.Dependency)
	}

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

func (worker *worker) handleStart(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error

	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = worker.handler.Start(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = worker.handler.StartVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStart, startStop)
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (worker *worker) handleStop(startStopMessage *messages.OrStartStop) (messages.NFVMessage, *vnfmError) {
	vnfr := startStopMessage.VNFR
	nsrID := vnfr.ParentNsID
	vnfcInstance := startStopMessage.VNFCInstance

	startStop := &messages.VNFMStartStop{VNFCInstance: vnfcInstance}

	var err error
	if vnfcInstance == nil { // Start the VNF Record
		startStop.VNFR, err = worker.handler.Stop(vnfr)
	} else { // Start the VNFC Instance
		startStop.VNFR, err = worker.handler.StopVNFCInstance(vnfr, vnfcInstance)
	}

	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionStop, startStop)
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (worker *worker) handleUpdate(updateMessage *messages.OrUpdate) (messages.NFVMessage, *vnfmError) {
	vnfr := updateMessage.VNFR
	nsrID := vnfr.ParentNsID
	script := updateMessage.Script

	replyVNFR, err := worker.handler.UpdateSoftware(script, vnfr)
	if err != nil {
		return nil, &vnfmError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.New(catalogue.ActionUpdate, &messages.VNFMGeneric{
		VNFR: replyVNFR,
	})
	if err != nil {
		worker.l.Panic("BUG: shouldn't happen")
	}

	return nfvMessage, nil
}

func (worker *worker) allocateResources(
	vnfr *catalogue.VirtualNetworkFunctionRecord,
	vimInstances map[string]interface{},
	keyPairs []*catalogue.Key) (*catalogue.VirtualNetworkFunctionRecord, *vnfmError) {

	worker.l.Debug("allocating resources for the VNFR")

	userData := worker.handler.UserData()

	worker.l.Debug("will send to NFVO UserData")

	msg, err := messages.New(&messages.VNFMAllocateResources{
		VNFR:         vnfr,
		VIMInstances: vimInstances,
		Userdata:     userData,
		KeyPairs:     keyPairs,
	})
	if err != nil {
		worker.l.Panicf("BUG")
	}

	nfvoResp, err := worker.executeRpc("vnfm.nfvo.actions.reply", msg)
	if err != nil {
		worker.l.Error("exchange error")

		return nil, &vnfmError{
			msg:   "Unable to allocate Resources",
			nsrID: vnfr.ParentNsID,
			vnfr:  vnfr,
		}
	}

	if nfvoResp != nil {
		if nfvoResp.Action() == catalogue.ActionError {
			errorMessage := nfvoResp.Content().(*messages.OrError)

			worker.l.Error("received error message from the NFVO")

			errVNFR := errorMessage.VNFR

			return nil, &vnfmError{
				msg:   fmt.Sprintf("Unable to allocate Resources. Reason: %s", errorMessage.Message),
				vnfr:  errVNFR,
				nsrID: vnfr.ParentNsID,
			}
		}

		message := nfvoResp.Content().(*messages.OrGeneric)
		worker.l.Debug("received a VNFR from ALLOCATE")

		return message.VNFR, nil
	}

	return nil, &vnfmError{
		msg:   "received an empty message from the NFVO",
		nsrID: vnfr.ParentNsID,
		vnfr:  vnfr,
	}
}
