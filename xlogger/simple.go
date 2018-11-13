package xlogger

import "log"

func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

type SimpleLogger struct {
}

func (*SimpleLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (*SimpleLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

func (*SimpleLogger) Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

func (*SimpleLogger) Debugf(format string, v ...interface{}) {
	log.Printf("[DEBUG] "+format, v...)
}

func (*SimpleLogger) Warnf(format string, v ...interface{}) {
	log.Printf("[WARN]  "+format, v...)
}

func (*SimpleLogger) Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
