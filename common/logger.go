package common

import (
	log "github.com/Sirupsen/logrus"
)

type Log interface {
	Info()
	Debug()
}

type LogClient struct {
	Location string
	Level    string
}

type LogServer struct {
	Location string
	Level    string
}

var logc = log.New()
var logs = log.New()

func (c LogClient) Info(args ...interface{}) {
	logc.Info(args...)
}

func (c LogClient) Debug(args ...interface{}) {
	logc.Debug(args...)
}

func (c LogClient) Warning(args ...interface{}) {
	logc.Warning(args...)
}

func (c LogClient) Error(args ...interface{}) {
	logc.Error(args...)
}

func (c LogClient) Panic(args ...interface{}) {
	logc.Panic(args...)
}

func (s LogServer) Info(args ...interface{}) {
	logs.Info(args...)
}

func (s LogServer) Debug(args ...interface{}) {
	logs.Debug(args...)
}

func (s LogServer) Warning(args ...interface{}) {
	logs.Warning(args...)
}

func (s LogServer) Error(args ...interface{}) {
	logs.Error(args...)
}

func (s LogServer) Panic(args ...interface{}) {
	logs.Panic(args...)
}
