package env

import "errors"

var (
	ErrVersionCompare  = errors.New("version compare error")
	ErrMissReqParam    = errors.New("missing required parameters")
	ErrSignature       = errors.New("signature error")
	ErrNeedLogin       = errors.New("please login first")
	ErrOperationFailed = errors.New("the operation failed. Please try again later")
)
