package messages


type NFVMessage interface {
	Action() Action
}