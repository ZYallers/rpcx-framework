package env

import "errors"

var (
	ErrVersionCompare       = errors.New("version compare error")
	ErrMissReqParam         = errors.New("missing required parameters")
	ErrSignature            = errors.New("signature error")
	ErrNeedLogin            = errors.New("please login first")
	ErrRegisterServiceEmpty = errors.New("after parsing, the service to be registered is empty")
	ErrOperationFailed = errors.New("the operation failed. Please try again later")
)
