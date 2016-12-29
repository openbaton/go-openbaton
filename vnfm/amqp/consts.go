package amqp

const (
	ExchangeDefault = "openbaton-exchange"
)

const (
	QueueVNFMRegister         = "nfvo.vnfm.register"
	QueueVNFMUnregister       = "nfvo.vnfm.unregister"
	QueueVNFMCoreActions      = "vnfm.nfvo.actions"
	QueueVNFMCoreActionsReply = "vnfm.nfvo.actions.reply"
	QueueNFVOGenericActions   = "nfvo.type.actions"
	QueueEMSRegistrator       = "ems.generic.register"
	QueueLogDispatch          = "nfvo.vnfm.logs"
)
