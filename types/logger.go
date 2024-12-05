package types

import (
	"fmt"
	"runtime/debug"

	libLogger "github.com/ZYallers/golib/utils/logger"
	"go.uber.org/zap"
)

type logger struct {
	Sender
	handler func() *zap.Logger
}

func NewLogger(name string, sender Sender) *logger {
	return &logger{
		Sender: sender,
		handler: func() *zap.Logger {
			return libLogger.Use(name)
		},
	}
}

func (l *logger) Debug(v ...interface{}) {
	l.handler().Debug(fmt.Sprint(v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.handler().Debug(fmt.Sprintf(format, v...))
}

func (l *logger) Info(v ...interface{}) {
	l.handler().Info(fmt.Sprint(v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	if format == "client has closed this connection: %s" {
		return
	}
	l.handler().Info(fmt.Sprintf(format, v...))
}

func (l *logger) Warn(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Warn(s)
	if l.Sender != nil {
		l.Sender.Error("Warn: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Warn(s)
	if l.Sender != nil {
		l.Sender.Error("Warn: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Error(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Error(s)
	if l.Sender != nil {
		l.Sender.Error("Error: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Error(s)
	if l.Sender != nil {
		l.Sender.Error("Error: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Fatal(s)
	if l.Sender != nil {
		l.Sender.Error("Fatal: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Fatal(s)
	if l.Sender != nil {
		l.Sender.Error("Fatal: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Panic(s)
	if l.Sender != nil {
		l.Sender.Error("Panic: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Panic(s)
	if l.Sender != nil {
		l.Sender.Error("Panic: "+s, string(debug.Stack()), true)
	}
}

func (l *logger) Handle(v ...interface{}) {
	l.Error(v...)
}
