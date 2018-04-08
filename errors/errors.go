package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var PropertyNotExists = errors.New("property not exists")

func ArgIsNil(argName string) error {
	return NewErrorFromString("argument %v is nil", argName)
}

func AssertTrue(name string, val bool, cause string) {
	if val == false {
		panic(NewErrorFromString("%s %s", name, cause))
	}
}

func AssertNotNil(name string, obj interface{}) {
	if obj == nil {
		panic(NewErrorFromString("%s is nil", name))
	}
}

func MustBeNotNil(argName string, obj interface{}) error {
	if obj == nil {
		return ArgIsNil(argName)
	}

	return nil
}

func NewErrorFromString(str string, formatArgs ...interface{}) error {
	return errors.New(fmt.Sprintf(str, formatArgs...))
}

func Wrap(e error, caller string) error {
	return errors.Wrap(e, caller)
}
