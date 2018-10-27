package log

import (
	"fmt"
	"github.com/go-stack/stack"
	log "github.com/inconshreveable/log15"
	"os"
)

type LogLevel string

const (
	DebugLevel    LogLevel = "debug"
	InfoLevel              = "info"
	ErrorLevel             = "error"
	CriticalLevel          = "critical"

	skipCallStackSteps = 8
)

type (
	logger struct{}
)

func CallerFileHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(r *log.Record) error {
		r.Ctx = append(r.Ctx, "caller", fmt.Sprintf("%+v", stack.Caller(skipCallStackSteps)))
		return h.Log(r)
	})
}

func Init(l LogLevel) {
	var logLvl log.Lvl

	switch l {
	case DebugLevel:
		logLvl = log.LvlDebug
	case InfoLevel:
		logLvl = log.LvlInfo
	case ErrorLevel:
		logLvl = log.LvlWarn
	case CriticalLevel:
		logLvl = log.LvlCrit
	default:
		fmt.Println("Logger: Uknown debug level. Defaulting to 'debug'")
		logLvl = log.LvlDebug
	}

	log.Root().SetHandler(log.LvlFilterHandler(logLvl, CallerFileHandler(log.StdoutHandler)))
}

func Debugf(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...interface{}) {
	log.Info(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...interface{}) {
	log.Warn(fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	log.Error(fmt.Sprintf(format, args...))
}

func Fatalf(format string, args ...interface{}) {
	log.Crit(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Debug(msg string, args ...interface{}) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	log.Crit(msg, args...)
	os.Exit(1)
}

func GetInstance() *logger {
	return &logger{}
}

func (this *logger) Debug(msg string, args ...interface{}) {
	log.Debug(msg, args...)
}

func (this *logger) Info(msg string, args ...interface{}) {
	log.Info(msg, args...)
}

func (this *logger) Warn(msg string, args ...interface{}) {
	log.Warn(msg, args...)
}

func (this *logger) Error(msg string, args ...interface{}) {
	log.Error(msg, args...)
}

func (this *logger) Fatal(msg string, args ...interface{}) {
	log.Crit(msg, args...)
	os.Exit(1)
}
