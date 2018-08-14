package vnfmsdk

import (
	"encoding/json"
	"fmt"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/streadway/amqp"
)

//Handler function for the VNFMs to be passed to the sdk package
func handleNfvMessage(bytemsg []byte, handlerVnfm sdk.Handler, allocate bool, connection *amqp.Connection, net catalogue.BaseNetworkInt, img catalogue.BaseImageInt) ([]byte, error) {
	logger := sdk.GetLogger("handler-function", "DEBUG")
	n, err := messages.Unmarshal(bytemsg, messages.NFVO)
	if err != nil {
		logger.Errorf("Error while unmarshaling nfv message: %v", err)
		err := sdk.NewSdkError(fmt.Sprintf("Error while unmarshaling nfv message: %v", err))
		return nil, err
	}
	logger.Debugf("Received Message %s", n.Action())
	switch h := handlerVnfm.(type) {
	case HandlerVnfm:
		wk := &worker{
			l:          logger,
			handler:    h,
			Allocate:   allocate,
			Connection: connection,
		}
		response := handleMessage(n, wk)
		var byteRes []byte
		resp, err := json.Marshal(response)
		if err != nil {
			logger.Errorf("Error while marshaling response: %v", err)
			return nil, err
		}
		byteRes = []byte(resp)
		return byteRes, nil
	default:
		return nil, sdk.NewSdkError("Not a HandlerVnfm implementation")
	}
}

func handleMessage(nfvMessage messages.NFVMessage, worker *worker) messages.NFVMessage {
	content := nfvMessage.Content()

	var reply messages.NFVMessage
	var err *vnfmError

	switch nfvMessage.Action() {

	case catalogue.ActionConfigure:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = worker.handleConfigure(genericMessage)

	case catalogue.ActionError:
		errorMessage := content.(*messages.OrError)
		err = worker.handleError(errorMessage)

	// case catalogue.ActionHeal:
	// 	healMessage := content.(*messages.OrHealVNFRequest)
	// 	reply, err = worker.handleHeal(healMessage)

	case catalogue.ActionInstantiate:
		instantiateMessage := content.(*messages.OrInstantiate)
		reply, err = worker.handleInstantiate(instantiateMessage)

	case catalogue.ActionInstantiateFinish:

	case catalogue.ActionModify:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = worker.handleModify(genericMessage)

	case catalogue.ActionScaleIn:
		scalingMessage := content.(*messages.OrScaling)
		err = worker.handleScaleIn(scalingMessage)

	case catalogue.ActionScaleOut:
		scalingMessage := content.(*messages.OrScaling)
		reply, err = worker.handleScaleOut(scalingMessage)

		// not implemented
	case catalogue.ActionScaling:

	// case catalogue.ActionStart:
	// 	startStopMessage := content.(*messages.OrStartStop)
	// 	reply, err = worker.handleStart(startStopMessage)

	// case catalogue.ActionStop:
	// 	startStopMessage := content.(*messages.OrStartStop)
	// 	reply, err = worker.handleStop(startStopMessage)

	// not implemented
	case catalogue.ActionReleaseResourcesFinish:

	case catalogue.ActionReleaseResources:
		genericMessage := content.(*messages.OrGeneric)
		reply, err = worker.handleReleaseResources(genericMessage)

	// case catalogue.ActionResume:
	// 	genericMessage := content.(*messages.OrGeneric)
	// 	reply, err = worker.handleResume(genericMessage)

	// case catalogue.ActionUpdate:
	// 	updateMessage := content.(*messages.OrUpdate)
	// 	reply, err = worker.handleUpdate(updateMessage)

	default:
		worker.l.Warning("received unsupported action")
	}
	if err != nil {
		worker.l.Errorf("%v", err)
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
			worker.l.Errorf("Error generating vnfm message error: %v", err)
		}
	}

	return reply
}
