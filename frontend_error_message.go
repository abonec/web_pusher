package web_pusher

import "encoding/json"

type FrontendErrorMessage struct {
	MsgType string            `json:"type"`
	Code    FrontendErrorCode `json:"code"`
	Message string            `json:"message"`
}

func (msg *FrontendErrorMessage) toJSON() ([]byte, error) {
	result, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}

	return result, nil
}

type FrontendErrorCode int

const (
	AuthError FrontendErrorCode = 1
)

func NewFrontendErrorMessage(code FrontendErrorCode, err error) *FrontendErrorMessage {
	return &FrontendErrorMessage{"error", code, err.Error()}
}
