package iferr_test

import (
	goerr "errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/unchainio/pkg/xlogger"

	"github.com/unchainio/pkg/iferr"

	"github.com/unchainio/pkg/errors"
)

func TestMessage(t *testing.T) {
	err := goerr.New("asdf")
	err = errors.Wrap(err, "qwe qwe1")
	err = errors.WithMessage(err, "qwe qwe2")
	err = errors.WithMessage(err, "qwe qwe3")
	err = errors.WithMessage(err, "qwe qwe4")

	err2 := goerr.New("asdf1")
	err2 = goerr.New("asdf2: " + err2.Error())
	err2 = goerr.New("asdf3: " + err2.Error())
	err2 = goerr.New("asdf4: " + err2.Error())

	fmt.Printf("%q\n", err)
	fmt.Printf("%q\n", err2)
}

func TestMessage2(t *testing.T) {
	fmt.Printf("%q", "as\"d\n\n")
}

func Test1(t *testing.T) {
	tcs := map[string]struct {
		opts     *iferr.ResponseOpts
		expected string
	}{
		"test1": {
			&iferr.ResponseOpts{},
			"error: err5: err4: err3: err2: err1\n",
		},
		"test2": {
			&iferr.ResponseOpts{Message: "qwe"},
			"qwe\n",
		},
		"test3": {
			&iferr.ResponseOpts{Depth: 1},
			"error: err5\n",
		},
		"test4": {
			&iferr.ResponseOpts{Depth: 2},
			"error: err5: err4\n",
		},
	}

	iferr := iferr.New(&xlogger.Mock{})

	for _, tc := range tcs {
		w := httptest.NewRecorder()

		iferr.Respond(w,
			errors.WithMessage(
				errors.WithMessage(
					errors.WithMessage(
						errors.Wrap(
							goerr.New("err1"),
							"err2",
						),
						"err3",
					),
					"err4",
				), "err5",
			),
			tc.opts,
		)

		require.Equal(t, tc.expected, w.Body.String())
	}

}
