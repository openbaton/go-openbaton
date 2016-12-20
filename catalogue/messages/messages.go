package messages

import "github.com/mcilloni/go-openbaton/catalogue"

type Content interface{}

type NFVMessage interface {
	Action() catalogue.Action
	Content() Content
	From() SenderType
}

// SenderType represents the type of the sender of the
// given message
type SenderType int

const (
	VNFM SenderType = iota
	NFVO
)

type message struct {
	action  catalogue.Action
	content Content
	from    SenderType
}

func (msg *message) Action() catalogue.Action {
	return msg.action
}

func (msg *message) Content() Content {
	return msg.content
}

func (msg *message) From() SenderType {
	return msg.from
}
