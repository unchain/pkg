package iferr

import (
	"fmt"
	"net/http"
	"os"

	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/xlogger"
)

var Default *IfErr

func init() {
	var err error

	log, err := xlogger.New(nil)

	if err != nil {
		panic(err)
	}

	Default = New(log)
}

type IfErr struct {
	log logger.Logger
}

func New(log logger.Logger) *IfErr {
	return &IfErr{
		log: log,
	}
}

func Warn(err error) { Default.Warn(err) }
func (ie *IfErr) Warn(err error) {
	if err != nil {
		ie.log.Warnf("%+v\n", err)
	}
}

func Exit(err error) { Default.Exit(err) }
func (ie *IfErr) Exit(err error) {
	if err != nil {
		ie.log.Errorf("%+v\n", err)
		os.Exit(-1)
	}
}

func Panic(err error) { Default.Panic(err) }
func (ie *IfErr) Panic(err error) {
	if err != nil {
		message := fmt.Sprintf("%+v\n", err)
		panic(message)
	}
}

type Fataler interface {
	Fatalf(format string, args ...interface{})
}

func Fail(f Fataler, err error) { Default.Fail(f, err) }
func (ie *IfErr) Fail(f Fataler, err error) {
	if err != nil {
		f.Fatalf("%+v\n", err)
	}
}

func WriteHTTP(err error, w http.ResponseWriter, codeArg ...int) bool {
	return Default.WriteHTTP(err, w, codeArg...)
}
func (ie *IfErr) WriteHTTP(err error, w http.ResponseWriter, codeArg ...int) bool {
	if err != nil {
		code := http.StatusInternalServerError

		ie.log.Errorf("HTTP Error: %+v\n", err)

		if len(codeArg) != 0 {
			code = codeArg[0]
		}

		http.Error(w, http.StatusText(code), code)

		return true
	}

	return false
}
