package xlogger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/unchainio/pkg/xmerge"
)

type Config struct {
	Level  string
	Format string
}

type Logger struct {
	log    *logrus.Logger
	errLog *logrus.Logger
	cfg    *Config
}

var defaultConfig = &Config{
	Level:  "debug",
	Format: "text",
}

func New(input *Config) (*Logger, error) {
	cfg := defaultConfig

	err := xmerge.Merge(cfg, input)

	if err != nil {
		return nil, errors.Wrap(err, "failed to merge with default config.")
	}

	logLevel, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		return nil, errors.Wrap(err, "unsupported log level")
	}

	var formatter logrus.Formatter

	switch cfg.Format {
	case "text":
		formatter = &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		}
	case "json":
		fallthrough
	default:
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.999Z07:00",
		}
	}

	l := &Logger{
		log: &logrus.Logger{
			Out:       os.Stdout,
			Hooks:     make(logrus.LevelHooks),
			Formatter: formatter,
			Level:     logLevel,
		},
		errLog: &logrus.Logger{
			Out:       os.Stderr,
			Hooks:     make(logrus.LevelHooks),
			Formatter: formatter,
			Level:     logLevel,
		},
		cfg: cfg,
	}

	return l, nil
}

func (l *Logger) Log(log *logrus.Logger) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(2); ok {
		fName := runtime.FuncForPC(pc).Name()

		currentDir, _ := os.Getwd()
		file = strings.Replace(file, currentDir, "", -1)

		caller := fmt.Sprintf("%s:%v", file, line)

		return log.WithField("caller", caller).WithField("fName", fName)
	}
	return &logrus.Entry{}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Log(l.log).Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log(l.errLog).Fatalf(format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Log(l.errLog).Panicf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Log(l.log).Debugf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Log(l.log).Warnf(format, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Log(l.log).Warningf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Log(l.errLog).Errorf(format, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.Log(l.log).Print(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Log(l.errLog).Fatal(v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.Log(l.errLog).Panic(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Log(l.log).Debug(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Log(l.log).Warn(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Log(l.log).Warning(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Log(l.errLog).Error(v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.Log(l.log).Println(v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Log(l.errLog).Fatalln(v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.Log(l.errLog).Panicln(v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	l.Log(l.log).Debugln(v...)
}

func (l *Logger) Warnln(v ...interface{}) {
	l.Log(l.log).Warnln(v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Log(l.log).Warningln(v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Log(l.errLog).Errorln(v...)
}
