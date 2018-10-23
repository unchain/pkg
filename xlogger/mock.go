package xlogger

type Mock struct {
}

func (*Mock) Printf(format string, v ...interface{}) {
	return
}

func (*Mock) Fatalf(format string, v ...interface{}) {
	return
}

func (*Mock) Panicf(format string, v ...interface{}) {
	return
}

func (*Mock) Debugf(format string, v ...interface{}) {
	return
}

func (*Mock) Warnf(format string, v ...interface{}) {
	return
}

func (*Mock) Errorf(format string, v ...interface{}) {
	return
}
