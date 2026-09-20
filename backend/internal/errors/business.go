package errors

type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string { return e.Message }

// Of 按错误码与消息构造业务异常，错误信息集中在调用处按场景给出。
func Of(code, message string) BusinessError {
	return BusinessError{Code: code, Message: message}
}

// Payload 参数类错误的快捷构造。
func Payload(message string) BusinessError {
	return Of("INVALID_PAYLOAD", message)
}
