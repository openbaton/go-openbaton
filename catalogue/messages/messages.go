package messages

import "github.com/mcilloni/go-openbaton/catalogue"
import "fmt"

type NFVMessage interface {
	Action() catalogue.Action
	Content() interface{}
	From() SenderType
}

// New creates dynamically a new NFVMessage from the given parameters.
// The function accepts a message body and an optional catalogue.Action before, such as in
// messages.NewMessage(catalogue.ActionError, &VNFMError{}).
// If no Action is specified, it is inferred using the DefaultAction() method of the body.
func New(params ...interface{}) (NFVMessage, error) {
	action := catalogue.NoActionSpecified
	var content body

	switch len(params) {
	case 1:
		if castContent, ok := params[0].(body); ok {
			content = castContent
		} else {
			return nil, fmt.Errorf("got %T, expected a valid message body type", params[0])
		}

		action = content.DefaultAction()


	case 2:
		if castAction, ok := params[0].(catalogue.Action); ok {
			action = castAction
		} else {
			return nil, fmt.Errorf("got %T, expected catalogue.Action", params[0])
		}

		if castContent, ok := params[1].(body); ok {
			content = castContent
		} else {
			return nil, fmt.Errorf("got %T, expected a valid message body type", params[1])
		}


	default:
		return nil, fmt.Errorf("wrong number of parameters for NewMessage(): %d", len(params))
	}

	return &message{
		action:  action,
		content: content,
	}, nil
}

// SenderType represents the type of the sender of the
// given message
type SenderType int

const (
	VNFM SenderType = iota
	NFVO
)

type body interface {
	DefaultAction() catalogue.Action
	From() SenderType
}

type message struct {
	action  catalogue.Action
	content body
}

func (msg *message) Action() catalogue.Action {
	return msg.action
}

func (msg *message) Content() interface{} {
	return msg.content
}

func (msg *message) From() SenderType {
	return msg.content.From()
}
