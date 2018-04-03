package utils

import (
	"../../log"
	"fmt"
	"github.com/pkg/errors"
)

func FatalExit(err error) {
	errStr := fmt.Sprintf("%+v", err)

	log.Global().Debug(`error details:
`+errStr, nil)
	log.Global().Fatal(err.Error()+`

exit app
`, nil)
	//log.Fatal(err)
	//panic(err)
}

func AssertExists(obj interface{}, name string) {
	if obj == nil {
		FatalExit(errors.New(fmt.Sprintf("%v not created", name)))
	}
}
