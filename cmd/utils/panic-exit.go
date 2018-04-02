package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"os"

	"../../prompt"
)

func PanicExit(err error) {
	prompt.Shared().Panic(err.Error())
	prompt.Shared().Debug("error details: %+v", err)
	//log.Fatal(err)
	//panic(err)
	os.Exit(1)
}

func AssertExists(obj interface{}, name string) {
	if obj == nil {
		PanicExit(errors.New(fmt.Sprintf("%v not created", name)))
	}
}
