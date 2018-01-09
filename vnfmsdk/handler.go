package vnfmsdk

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
)

func handleNfvMessage(bytemsg []byte, wk interface{}) ([]byte, error) {
	logger := sdk.GetLogger("handler-function", "DEBUG")
	n, verr := messages.Unmarshal(bytemsg, messages.NFVO)
	if verr != nil {
		logger.Errorf("Error while unmarshaling nfv message: %v", verr)
		err := errors.New(fmt.Sprintf("Error while unmarshaling nfv message: %v", verr))
		return nil, err
	}
	logger.Debugf("Received Message %s", n.Action())
	var response messages.NFVMessage
	switch t := wk.(type) {

	case *worker:
		response = handleMessage(n, t)
	default:
		return nil, errors.New("Worker must be of type ")
	}
	var byteRes []byte
	resp, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Error while marshaling response: %v", err)
		return nil, err
	}
	byteRes = []byte(resp)
	return byteRes, nil
}

func handleMessage(nfvMessage messages.NFVMessage, wk *worker) (messages.NFVMessage) {
	content := nfvMessage.Content()

	var reply messages.NFVMessage
	var err *vnfmError

	switch nfvMessage.Action() {

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
		wk.l.Warning("received unsupported action")
	}
	if err != nil {
		wk.l.Errorf("%v", err)
		errorMsg, err := messages.New(catalogue.ActionError, &messages.VNFMError{
			Exception: messages.JavaException{
				DetailMessage:        err.msg,
				StackTrace:           make([]messages.Trace, 0),
				SuppressedExceptions: make([]string, 0),
				InternalCause: messages.Cause{
					DetailMessage:        err.msg,
					StackTrace:           make([]messages.Trace, 0),
					SuppressedExceptions: make([]string, 0),
				},
			},
			NSRID: err.nsrID,
			VNFR:  err.vnfr,
		})
		if err == nil {
			return errorMsg
		} else {
			wk.l.Errorf("Error generating vnfm message error: %v", err)
		}
	}

	return reply
}
