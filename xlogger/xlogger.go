package xlogger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type Config struct {
	Level  string
	Format string
}

type Logger struct {
	log *logrus.Logger
	cfg *Config
}

func New(cfg *Config) (*Logger, error) {

	if cfg.Level == "" {
		cfg.Level = "debug"
	}

	logLevel, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		return nil, errors.Wrap(err, "Unsupported log level")
	}

	var formatter logrus.Formatter

	switch cfg.Format {
	case "text":
		formatter = &logrus.TextFormatter{
		//DisableColors: true,
		//DisableSorting:   true,
		//QuoteEmptyFields: false,
		}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.999Z07:00",
		}
	}

	l := &Logger{
		cfg: cfg,
		log: &logrus.Logger{
			Out:       os.Stderr,
			Formatter: formatter,
			Hooks:     make(logrus.LevelHooks),
			Level:     logLevel,
		},
	}

	return l, nil
}

func (l *Logger) Log() *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(2); ok {
		fName := runtime.FuncForPC(pc).Name()

		currentDir, _ := os.Getwd()
		file = strings.Replace(file, currentDir, "", -1)

		caller := fmt.Sprintf("%s:%v", file, line)

		return l.log.WithField("caller", caller).WithField("fName", fName)
	}
	return &logrus.Entry{}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Log().Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log().Fatalf(format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Log().Panicf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Log().Debugf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Log().Warnf(format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Log().Warningf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Log().Errorf(format, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.Log().Print(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Log().Fatal(v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.Log().Panic(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Log().Debug(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Log().Warn(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Log().Warning(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Log().Error(v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.Log().Println(v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Log().Fatalln(v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.Log().Panicln(v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	l.Log().Debugln(v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	l.Log().Warnln(v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Log().Warningln(v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Log().Errorln(v...)
}
