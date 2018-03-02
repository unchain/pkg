package xio

import (
	"fmt"
	"io"
)

// MultiCopy copies a reader to multiple writers
func MultiCopy(src io.Reader, dsts ...io.Writer) {
	var err error

	buff := make([]byte, 100)

	var n int
	for err == nil {
		n, err = src.Read(buff)
		if n > 0 {
			for _, dst := range dsts {
				fmt.Fprintf(dst, "%s", string(buff[:n]))
			}
		}
	}
}
