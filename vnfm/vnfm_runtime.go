package vnfm

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	log "github.com/sirupsen/logrus"
	"github.com/mcilloni/go-openbaton/catalogue"
	"strings"
)

type logger *log.Logger

type VNFM struct {
	logger
	conn     NFVOConnector
	impl      Provider
	props Properties
	quitChan chan struct{}
}

func Start(connector NFVOConnector, impl Provider, properties Properties) (*VNFM, error) {
	vnfm := &VNFM{
		conn:     connector,
		impl:      impl,
		logger:   logger(log.New()),
		props: properties,
		quitChan: make(chan struct{}),
	}

	msgChan, err := connector.NotifyReceived()
	if err != nil {
		return nil, err
	}

	go vnfm.loop(msgChan)

	return vnfm, nil
}

func (vnfm *VNFM) Logger() *log.Logger {
	return vnfm.logger
}

func (vnfm *VNFM) SetLogger(log *log.Logger) {
	vnfm.logger = (*logger)(log)
}

func (vnfm *VNFM) Stop() error {
	vnfm.quitChan <- struct{}{}

	select {
	case <-vnfm.quitChan:
		return nil
	case <-time.After(time.Second):
		return errors.New("the VNFM refused to quit")
	}
}

type vnfrError struct {
	msg string
	vnfr *catalogue.VirtualNetworkFunctionRecord
	nsrID string
}

func (e *vnfrError) Error() string {
	return msg
}

func (vnfm *VNFM) handle(message messages.NFVMessage) {

	vnfm.Debugln("VNFM: Received Message: " + message.Action())

	content := message.Content()

	var reply messages.NFVMessage
	var err *vnfrError

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

	case RELEASE_RESOURCES_FINISH:
		break
	case UPDATE:
		OrVnfmUpdateMessage orVnfmUpdateMessage = (OrVnfmUpdateMessage) message
		nfvMessage =
			VnfmUtils.getNfvMessage(
				Action.UPDATE,
				updateSoftware(orVnfmUpdateMessage.getScript(), orVnfmUpdateMessage.getVnfr()))
		break
	case HEAL:
		OrVnfmHealVNFRequestMessage orVnfmHealMessage = (OrVnfmHealVNFRequestMessage) message

		nsrId = orVnfmHealMessage.getVirtualNetworkFunctionRecord().getParent_ns_id()
		VirtualNetworkFunctionRecord vnfrObtained =
			this.heal(
				orVnfmHealMessage.getVirtualNetworkFunctionRecord(),
				orVnfmHealMessage.getVnfcInstance(),
				orVnfmHealMessage.getCause())
		nfvMessage =
			VnfmUtils.getNfvMessageHealed(
				Action.HEAL, vnfrObtained, orVnfmHealMessage.getVnfcInstance())

		break
	case INSTANTIATE_FINISH:
		break
	case CONFIGURE:
		orVnfmGenericMessage = (OrVnfmGenericMessage) message
		vnfr = orVnfmGenericMessage.getVnfr()
		nsrId = orVnfmGenericMessage.getVnfr().getParent_ns_id()
		nfvMessage =
			VnfmUtils.getNfvMessage(Action.CONFIGURE, configure(orVnfmGenericMessage.getVnfr()))
		break
	case START:
		{
		OrVnfmStartStopMessage orVnfmStartStopMessage = (OrVnfmStartStopMessage) message
		vnfr = orVnfmStartStopMessage.getVirtualNetworkFunctionRecord()
		nsrId = vnfr.getParent_ns_id()
		VNFCInstance vnfcInstance = orVnfmStartStopMessage.getVnfcInstance()

		if (vnfcInstance == null) // Start the VNF Record
		{
			nfvMessage =
				VnfmUtils.getNfvMessage(Action.START, start(vnfr))
		} else // Start the VNFC Instance
		{
			nfvMessage =
				VnfmUtils.getNfvMessageStartStop(
					Action.START,
					startVNFCInstance(vnfr, vnfcInstance),
					vnfcInstance)
		}
		break
		}
	case STOP:
		{
		OrVnfmStartStopMessage orVnfmStartStopMessage = (OrVnfmStartStopMessage) message
		vnfr = orVnfmStartStopMessage.getVirtualNetworkFunctionRecord()
		nsrId = vnfr.getParent_ns_id()
		VNFCInstance vnfcInstance = orVnfmStartStopMessage.getVnfcInstance()

		if (vnfcInstance == null) // Stop the VNF Record
		{
			nfvMessage = VnfmUtils.getNfvMessage(Action.STOP, stop(vnfr))
		} else // Stop the VNFC Instance
		{
			nfvMessage =
				VnfmUtils.getNfvMessageStartStop(
					Action.STOP,
					stopVNFCInstance(vnfr, vnfcInstance),
					vnfcInstance)
		}

		break
		}
	case RESUME:
		{
		OrVnfmGenericMessage orVnfmResumeMessage = (OrVnfmGenericMessage) message
		vnfr = orVnfmResumeMessage.getVnfr()
		nsrId = vnfr.getParent_ns_id()

		Action resumedAction = this.getResumedAction(vnfr, null)
		nfvMessage =
			VnfmUtils.getNfvMessage(
				resumedAction,
				resume(vnfr, null, orVnfmResumeMessage.getVnfrd()))
		log.debug(
			"Resuming vnfr '"
				+ vnfr.getId()
				+ "' with dependency target: '"
				+ orVnfmResumeMessage.getVnfrd().getTarget()
				+ "' for action: "
				+ resumedAction
				+ "'")
		break
		}
	}

	log.debug(
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	if (nfvMessage != null) {
	log.debug("send " + nfvMessage.getClass().getSimpleName() + " to NFVO")
	vnfmHelper.sendToNfvo(nfvMessage)
	}
	
	catch (Throwable e) {
      log.error("ERROR: ", e)
      if (e instanceof VnfmSdkException) {
        VnfmSdkException vnfmSdkException = (VnfmSdkException) e
        if (vnfmSdkException.getVnfr() != null) {
          log.debug("sending VNFR with version: " + vnfmSdkException.getVnfr().getHb_version())
          vnfmHelper.sendToNfvo(
              VnfmUtils.getNfvErrorMessage(vnfmSdkException.getVnfr(), vnfmSdkException, nsrId))
          return
        }
      } else if (e.getCause() instanceof VnfmSdkException) {
        VnfmSdkException vnfmSdkException = (VnfmSdkException) e.getCause()
        if (vnfmSdkException.getVnfr() != null) {
          log.debug("sending VNFR with version: " + vnfmSdkException.getVnfr().getHb_version())
          vnfmHelper.sendToNfvo(
              VnfmUtils.getNfvErrorMessage(vnfmSdkException.getVnfr(), vnfmSdkException, nsrId))
          return
        }
      }
      vnfmHelper.sendToNfvo(VnfmUtils.getNfvErrorMessage(vnfr, e, nsrId))
    }
}

func (vnfm *VNFM) handleError(errorMessage *messages.OrError) *vnfrError {
	vnfr := errorMessage.VNFR
	nsrID := vnfr.ParentNsID
	
	vnfm.Errorf("received an error from the NFVO: %s\n", errorMessage.Message)

	if err := vnfm.impl.HandleError(errorMessage.VNFR); err != nil {
		return &vnfrError{err.Error(), nsrID, vnfr}
	}

	return nil
}

func (vnfm *VNFM) handleInstantiate(instantiateMessage *messages.OrInstantiate) (messages.NFVMessage, *vnfrError) {
	extension := instantiateMessage.Extension

	vnfm.Debugf("received extensions: %v\n", extension);
	vnfm.Debugf("received keys: %v\n", instantiateMessage.Keys)

	vimInstances := instantiateMessage.VIMInstances
	
	vnfr =
		createVirtualNetworkFunctionRecord(
			orVnfmInstantiateMessage.getVnfd(),
			orVnfmInstantiateMessage.getVnfdf().getFlavour_key(),
			orVnfmInstantiateMessage.getVlrs(),
			orVnfmInstantiateMessage.getExtension(),
			vimInstances)
	GrantOperation grantOperation = new GrantOperation()
	grantOperation.setVirtualNetworkFunctionRecord(vnfr)

	Future<OrVnfmGrantLifecycleOperationMessage> result = executor.submit(grantOperation)
	OrVnfmGrantLifecycleOperationMessage msg
	try {
	msg = result.get()
	if (msg == null) {
		return
	}
	} catch (ExecutionException e) {
	log.error("Got exception while allocating vms")
	throw e.getCause()
	}

	vnfr = msg.getVirtualNetworkFunctionRecord()
	Map<String, VimInstance> vimInstanceChosen = msg.getVduVim()

	vnfm.Debugln("VERSION IS: " + vnfr.getHb_version())

	if (!properties.getProperty("allocate", "true").equalsIgnoreCase("true")) {
	AllocateResources allocateResources = new AllocateResources()
	allocateResources.setVirtualNetworkFunctionRecord(vnfr)
	allocateResources.setVimInstances(vimInstanceChosen)
	allocateResources.setKeyPairs(orVnfmInstantiateMessage.getKeys())
	try {
		vnfr = executor.submit(allocateResources).get()
		if (vnfr == null) {
		return
		}
	} catch (ExecutionException e) {
		log.error("Got exception while allocating vms")
		throw e.getCause()
	}
	}
	setupProvides(vnfr)

	for (VirtualDeploymentUnit virtualDeploymentUnit :
		vnfr.getVdu()) {
	for (VNFCInstance vnfcInstance : virtualDeploymentUnit.getVnfc_instance()) {
		checkEMS(vnfcInstance.getHostname())
	}
	}
	if (orVnfmInstantiateMessage.getVnfPackage() != null) {
	if (orVnfmInstantiateMessage.getVnfPackage().getScriptsLink() != null) {
		vnfr =
			instantiate(
				vnfr,
				orVnfmInstantiateMessage.getVnfPackage().getScriptsLink(),
				vimInstances)
	} else {
		vnfr =
			instantiate(
				vnfr,
				orVnfmInstantiateMessage.getVnfPackage().getScripts(),
				vimInstances)
	}
	} else {
	vnfr =
		instantiate(vnfr, null, vimInstances)
	}
	nfvMessage = VnfmUtils.getNfvMessage(Action.INSTANTIATE, vnfr)
	break
}

func (vnfm *VNFM) handleModify(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfrError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := vnfm.impl.Modify(vnfr, genericMessage.VNFRDependency)
	if err != nil {
		return nil, &vnfrError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.NewMessage(messages.ActionModify, &messages.VNFMGeneric{
		VNFR: resultVNFR,
	})
	if err != nil {
		return nil, &vnfrError{err.Error(), resultVNFR, nsrID}
	}

	return nfvMessage, nil
}

func (vnfm *VNFM) handleReleaseResources(genericMessage *messages.OrGeneric) (messages.NFVMessage, *vnfrError) {
	vnfr := genericMessage.VNFR
	nsrID := vnfr.ParentNsID

	resultVNFR, err := vnfm.impl.Terminate(vnfr)
	if err != nil {
		return nil, &vnfrError{err.Error(), vnfr, nsrID}
	}

	nfvMessage, err := messages.NewMessage(messages.ActionReleaseResources, &messages.VNFMGeneric{
		VNFR: resultVNFR,
	})
	if err != nil {
		return nil, &vnfrError{err.Error(), resultVNFR, nsrID}
	}

	return nfvMessage, nil
}

func (vnfm *VNFM) handleScaleIn(scalingMessage *messages.OrScaling) *vnfrError {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	
	vnfcInstanceToRemove := scalingMessage.getVnfcInstance()

	if _, err := vnfm.impl.Scale(catalogue.ActionScaleIn, vnfr, vnfcInstanceToRemove, nil, nil); err != nil {
		return &vnfrError{err.Error(), nsrID, vnfr}
	}
	
	return nil
}

func (vnfm *VNFM) handleScaleOut(scalingMessage *messages.OrScaling) (messages.NFVMessage, *vnfrError) {
	vnfr := scalingMessage.VNFR
	nsrID := vnfr.ParentNsID
	component := scalingMessage.Component

	vnfm.Debugf("HB_VERSION == %d\n", vnfr.HbVersion)
	vnfm.Infof("Adding VNFComponent: %v\n" + component)
	vnfm.Debugf("The mode is: %s\n", mode)

	var newVNFCInstance *catalogue.VNFCInstance
	if allocate, set := vnfm.props.ValueBool("vnfm.allocate"); set && allocate {
		newMsg, err := messages.NewMessage(&messages.VNFMScaling{
			VNFR: vnfr,
			UserData: vnfm.impl.UserData(),
		})

		if err != nil {
			return nil, &vnfrError{err.Error(), nsrID, vnfr}
		}

		var replyVNFR *catalogue.VirtualNetworkFunctionRecord

		switch content := newMsg.Content().(type) {
		case messages.OrGeneric:
			replyVNFR = content.VNFR
			vnfm.Debugf("HB_VERSION == %d\n", replyVNFR.HbVersion)
		
		case messages.OrError:
			if err := vnfm.impl.HandleError(content.VNFR); err != nil {
				return nil, &vnfrError{err.Error(), nsrID, content.VNFR}
			}

			return nil, nil

		default:
			replyVNFR = vnfr
		}

		if newVNFCInstance = replyVNFR.FindComponentInstance(component); newVNFCInstance == nil {
			return nil, vnfrError{"no new VNFCInstance found. This should not happen.", nsrID, replyVNFR}
		}

		vnfm.Debugf("VNFComponentInstance FOUND : %v\n", newVNFCInstance.VNFComponent)

		if strings.EqualFold(scalingMessage.Mode, "STANDBY") {
			newVNFCInstance.State = "STANDBY"
		}

		vnfm.impl.CheckEMS(newVNFCInstance.Hostname)

		vnfr = replyVNFR
	} else {
		vnfm.Warnln("vnfm.allocate is not set. No new VNFCInstance has been instantiated.")
	}

	var scripts interface{}

	switch {
	case scalingMessage.VNFPackage == nil:
		scripts = []*catalogue.Script{}

	case scalingMessage.VNFPackage.ScriptsLink != nil 
		scripts = scalingMessage.VNFPackage.ScriptsLink

	default:
		scripts = scalingMessage.VNFPackage.Scripts
	}

	resultVNFR, err := vnfm.impl.Scale(catalogue.ActionScaleOut, vnfr, newVNFCInstance, scripts, scalingMessage.Dependency)
	if err != nil {
		return nil, &vnfrError{err.Error(), nsrID, vnfr}
	}

	nfvMessage, err := messages.NewMessage(Action.SCALED, &messages.VNFMScaled{
		VNFR: resultVNFR,
		VNFCInstance: newVNFCInstance,
	})

	if err != nil {
		return nil, &vnfrError{err.Error(), nsrID, resultVNFR}
	}

	return nfvMessage, nil
}

func (vnfm *VNFM) loop(msgChan <-chan messages.NFVMessage) {
	defer func() {
		if r := recover(); r != nil {
			if err := vnfm.conn.Close(); err != nil {
				vnfm.Errorln(err)
			}
			vnfm.Logger.Panicln(r)
		}
	}()

MainLoop:
	for {
		select {
		case msg := <-msgChan:
			if err := vnfm.handle(msg); err != nil {
				vnfm.Errorln(err)
			}			

		case <-vnfm.quitChan:
			break MainLoop

		default:

		}
	}
}