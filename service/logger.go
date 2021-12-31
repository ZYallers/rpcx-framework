package service

import (
	"fmt"
	"github.com/ZYallers/rpcx-framework/helper"
	zap2 "github.com/ZYallers/rpcx-framework/util/zap"
	"go.uber.org/zap"
	"runtime/debug"
)

type logger struct {
	handler func() *zap.Logger
}

func NewLogger(name string) *logger {
	return &logger{handler: func() *zap.Logger {
		return zap2.Use(name)
	}}
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
	helper.ContextMessage("Warn: "+s, string(debug.Stack()), false)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Warn(s)
	helper.ContextMessage("Warn: "+s, string(debug.Stack()), false)
}

func (l *logger) Error(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Error(s)
	helper.ContextMessage("Error: "+s, string(debug.Stack()), true)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Error(s)
	helper.ContextMessage("Error: "+s, string(debug.Stack()), true)
}

func (l *logger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Fatal(s)
	helper.ContextMessage("Fatal: "+s, string(debug.Stack()), true)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Fatal(s)
	helper.ContextMessage("Fatal: "+s, string(debug.Stack()), true)
}

func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.handler().Panic(s)
	helper.ContextMessage("Panic: "+s, string(debug.Stack()), true)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.handler().Panic(s)
	helper.ContextMessage("Panic: "+s, string(debug.Stack()), true)
}

func (l *logger) Handle(v ...interface{}) {
	l.Error(v...)
}
