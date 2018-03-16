package xerrors

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
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
	return errors.WithStack(
		&warning{
			error: err,
		},
	)
}

type MultiError []error

func (err *MultiError) Error() string {

}
