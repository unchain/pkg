package xerrors

import (
	"fmt"
	"log"
	"os"
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
