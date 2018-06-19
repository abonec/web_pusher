package web_pusher

type FrontendMessage struct {
	MsgType    string            `json:"type"`
	AuthStatus string            `json:"auth_status"`
	Code       FrontendErrorCode `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	UserId     string            `json:"user_id,omitempty"`
}

//func (msg *FrontendMessage) toJSON() ([]byte, error) {
//	result, err := json.Marshal(msg)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return result, nil
//}

type FrontendErrorCode int

const (
	AuthError FrontendErrorCode = 1
)

func NewFrontendErrorMessage(code FrontendErrorCode, err error) *FrontendMessage {
	return &FrontendMessage{MsgType: "auth", AuthStatus: "failure", Code: code, Message: err.Error()}
}

func NewFrontendSuccessMessage(userId string) *FrontendMessage {
	return &FrontendMessage{MsgType: "auth", AuthStatus: "success", UserId: userId}
}
