package errors

type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string { return e.Message }

// 业务错误码
const (
	CodeCreditDenied       = "CREDIT_DENIED"
	CodeInvitationLimit    = "INVITATION_LIMIT"
	CodeInvitationFinal    = "INVITATION_ALREADY_FINAL"
	CodeInvitationNotFound = "INVITATION_NOT_FOUND"
	CodeNeedClosed         = "NEED_CLOSED"
	CodeDuplicateInvite    = "DUPLICATE_INVITATION"
)

func New(code, message string) BusinessError {
	return BusinessError{Code: code, Message: message}
}
