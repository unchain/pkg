package xerrors

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/pkg/errors"
	"github.com/unchainio/interfaces/logger"
)

func PanicIfError(err error) {
	if err != nil {
		message := fmt.Sprintf("%+v\n", err)
		panic(message)
	}
}

func ExitIfError(err error) {
	if err != nil {
		log.Printf("%+v\n", err)
		os.Exit(-1)
	}
}

func HTTPError(log logger.Logger, w http.ResponseWriter, err error) {
	if err != nil {
		log.Errorf("HTTP Error: %+v\n", err)

		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

type Warn interface {
	Warn() string
}

type Info interface {
	Info() string
}

type warning struct {
	error
}

func (w *warning) IsWarn() bool {
	return true
}

type warner interface {
	IsWarn() bool
}

func IsWarn(err error) bool {
	w, ok := errors.Cause(err).(warner)

	return ok && w.IsWarn()
}

func ToWarn(err error) error {
	return &warning{
		error: err,
	}
}
