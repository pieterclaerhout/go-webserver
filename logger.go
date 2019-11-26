package webserver

import (
	"encoding/json"
	"io"

	glog "github.com/labstack/gommon/log"

	"github.com/pieterclaerhout/go-log"
)

type logWriter struct{}

func (w logWriter) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), nil
}

type logger struct {
	w      io.Writer
	prefix string
	level  glog.Lvl
}

func newLogger() logger {
	return logger{
		w:      logWriter{},
		prefix: "",
		level:  glog.Lvl(2),
	}
}

func (l logger) Output() io.Writer {
	return l.w
}

func (l logger) SetOutput(w io.Writer) {
	l.w = w
}

func (l logger) Prefix() string {
	return l.prefix
}

func (l logger) SetPrefix(p string) {
	l.prefix = p
}

func (l logger) Level() glog.Lvl {
	return l.level
}

func (l logger) SetLevel(v glog.Lvl) {
	l.level = v
}

func (l logger) SetHeader(h string) {
	// Not implemented
}

func (l logger) Print(i ...interface{}) {
	log.Info(i...)
}

func (l logger) Printf(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (l logger) Printj(j glog.JSON) {
	log.Info(l.jsonToString(j))
}

func (l logger) Debug(i ...interface{}) {
	log.Debug(i...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func (l logger) Debugj(j glog.JSON) {
	log.Debug(l.jsonToString(j))
}

func (l logger) Info(i ...interface{}) {
	log.Info(i...)
}

func (l logger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func (l logger) Infoj(j glog.JSON) {
	log.Info(l.jsonToString(j))
}

func (l logger) Warn(i ...interface{}) {
	log.Warn(i...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func (l logger) Warnj(j glog.JSON) {
	log.Warn(l.jsonToString(j))
}

func (l logger) Error(i ...interface{}) {
	log.Error(i...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func (l logger) Errorj(j glog.JSON) {
	log.Error(l.jsonToString(j))
}

func (l logger) Fatal(i ...interface{}) {
	log.Fatal(i...)
}

func (l logger) Fatalj(j glog.JSON) {
	log.Fatal(l.jsonToString(j))
}

func (l logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func (l logger) Panic(i ...interface{}) {
	// Not implemented
}

func (l logger) Panicj(j glog.JSON) {
	// Not implemented
}

func (l logger) Panicf(format string, args ...interface{}) {
	// Not implemented
}

func (l logger) jsonToString(j glog.JSON) string {
	log.InfoDump(j, "")
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(b)
}
