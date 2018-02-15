package xmock

type Logger struct{}

func (Logger) Printf(format string, v ...interface{}) {}

func (Logger) Fatalf(format string, v ...interface{}) {}

func (Logger) Panicf(format string, v ...interface{}) {}

func (Logger) Debugf(format string, v ...interface{}) {}

func (Logger) Warnf(format string, v ...interface{}) {}

func (Logger) Warningf(format string, v ...interface{}) {}

func (Logger) Errorf(format string, v ...interface{}) {}

func (Logger) Print(v ...interface{}) {}

func (Logger) Fatal(v ...interface{}) {}

func (Logger) Panic(v ...interface{}) {}

func (Logger) Debug(v ...interface{}) {}

func (Logger) Warn(v ...interface{}) {}

func (Logger) Warning(v ...interface{}) {}

func (Logger) Error(v ...interface{}) {}

func (Logger) Println(v ...interface{}) {}

func (Logger) Fatalln(v ...interface{}) {}

func (Logger) Panicln(v ...interface{}) {}

func (Logger) Debugln(v ...interface{}) {}

func (Logger) Warnln(v ...interface{}) {}

func (Logger) Warningln(v ...interface{}) {}

func (Logger) Errorln(v ...interface{}) {}

