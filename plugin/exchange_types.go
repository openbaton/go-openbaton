package plugin

import (
	"encoding/json"
)

type plugError struct {
	Message string `json:"detailMessage"`
}

func (e plugError) Error() string {
	return e.Message
}

type request struct {
	MethodName string            `json:"methodName"`
	Parameters []json.RawMessage `json:"parameters"`

	// ReplyTo is the query onto which the reply should be sent
	ReplyTo string `json:"-"`

	// CorrID to be used while sending the reply
	CorrID string `json:"-"`
}

type response struct {
	Answer    interface{} `json:"answer,omitempty"`
	Exception error       `json:"exception,omitempty"`
}
