package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	NoGenesisBlock = errors.New("no genesis block")
)

func TxValidationFailed(reason string, formatArgs ...interface{}) error {
	if len(formatArgs) > 0 {
		reason = fmt.Sprintf(reason, formatArgs)
	}

	return NewErrorFromString("TxValidation FAILED: %s ", reason)
}

func TxValidationUnknownTransactionType(t string) error {
	return TxValidationFailed("unknown transaction type %s", t)
}

func ArgIsNil(argName string) error {
	return NewErrorFromString("argument %v is nil", argName)
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
