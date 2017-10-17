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
	var err *vnfmError
	switch t := wk.(type) {

	case *worker:
		response, err = handleMessage(n, t)
	default:
		return nil, errors.New("Worker must be of type ")
	}
	var byteRes []byte
	if err != nil {
		merr, err := json.Marshal(err)
		if err != nil {
			logger.Errorf("Error while marshaling error: %v", err)
			return nil, err
		}
		byteRes = []byte(merr)
	} else {
		resp, err := json.Marshal(response)
		//logger.Debugf("Sending back: \n%v" , string(resp))
		if err != nil {
			logger.Errorf("Error while marshaling response: %v", err)
			return nil, err
		}
		byteRes = []byte(resp)
	}
	return byteRes, nil
}

func handleMessage(nfvMessage messages.NFVMessage, wk *worker) (messages.NFVMessage, *vnfmError) {
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
		wk.l.Errorf("ERROR: %v", err)
	}

	return reply, err
}
