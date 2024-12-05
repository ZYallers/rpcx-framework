package sender

import (
	"github.com/ZYallers/rpcx-framework/types"
)

var sender types.Sender

func Register(s types.Sender) { sender = s }

func GetSender() types.Sender { return sender }

func Error(msg interface{}, stack string, isAtAll bool, logType ...interface{}) {
	if sender != nil {
		sender.Error(msg, stack, isAtAll, logType...)
	}
}

func Graceful(msg interface{}, isAtAll bool, logType ...interface{}) {
	if sender != nil {
		sender.Graceful(msg, isAtAll, logType...)
	}
}

func Push(token string, msg interface{}, options ...interface{}) {
	if sender != nil {
		sender.Send(token, msg, options...)
	}
}
