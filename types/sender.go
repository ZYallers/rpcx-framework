package types

type Sender interface {
	Error(msg interface{}, stack string, isAtAll bool, logType ...interface{})
	Graceful(msg interface{}, isAtAll bool, logType ...interface{})
	Send(token string, msg interface{}, options ...interface{})
}
