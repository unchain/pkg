package xhttp

import (
	"encoding/json"
	"io"

	"bitbucket.org/unchain/pkg/xerrors"
)

func WriteAsJSON(w io.Writer, obj interface{}) {
	bytes, err := json.Marshal(obj)
	xerrors.PanicIfError(err)

	w.Write(bytes)
}
