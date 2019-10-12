package webserver

import (
	"encoding/json"
	"io"

	glog "github.com/labstack/gommon/log"

	"github.com/pieterclaerhout/go-log"
)

type LogWriter struct{}

func (logWriter LogWriter) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), nil
}

type Logger struct {
	w      io.Writer
	prefix string
	level  glog.Lvl
}

func NewLogger() Logger {
	return Logger{
		w:      LogWriter{},
		prefix: "",
		level:  glog.Lvl(2),
	}
}

func (logger Logger) Output() io.Writer {
	return logger.w
}

func (logger Logger) SetOutput(w io.Writer) {
	logger.w = w
}

func (logger Logger) Prefix() string {
	return logger.prefix
}

func (logger Logger) SetPrefix(p string) {
	logger.prefix = p
}

func (logger Logger) Level() glog.Lvl {
	return logger.level
}

func (logger Logger) SetLevel(v glog.Lvl) {
	logger.level = v
}

func (logger Logger) SetHeader(h string) {
}

func (logger Logger) Print(i ...interface{}) {
	log.Info(i...)
}

func (logger Logger) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (logger Logger) Printj(j glog.JSON) {
	log.Info(logger.jsonToString(j))
}

func (logger Logger) Debug(i ...interface{}) {
	log.Debug(i...)
}

func (logger Logger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (logger Logger) Debugj(j glog.JSON) {
	log.Debug(logger.jsonToString(j))
}

func (logger Logger) Info(i ...interface{}) {
	log.Info(i...)
}

func (logger Logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (logger Logger) Infoj(j glog.JSON) {
	log.Info(logger.jsonToString(j))
}

func (logger Logger) Warn(i ...interface{}) {
	log.Warn(i...)
}

func (logger Logger) Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (logger Logger) Warnj(j glog.JSON) {
	log.Warn(logger.jsonToString(j))
}

func (logger Logger) Error(i ...interface{}) {
	log.Error(i...)
}

func (logger Logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (logger Logger) Errorj(j glog.JSON) {
	log.Error(logger.jsonToString(j))
}

func (logger Logger) Fatal(i ...interface{}) {
	log.Fatal(i...)
}

func (logger Logger) Fatalj(j glog.JSON) {
	log.Fatal(logger.jsonToString(j))
}

func (logger Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (logger Logger) Panic(i ...interface{}) {
}

func (logger Logger) Panicj(j glog.JSON) {
}

func (logger Logger) Panicf(format string, args ...interface{}) {
}

func (logger Logger) jsonToString(j glog.JSON) string {
	log.InfoDump(j, "")
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}
