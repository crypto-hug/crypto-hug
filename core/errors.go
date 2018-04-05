package core

import (
	"fmt"

	"github.com/crypto-hug/crypto-hug/errors"
)

var (
	NoGenesisBlock = errors.NewErrorFromString("no genesis block")
)

func TxValidationFailed(reason string, formatArgs ...interface{}) error {
	if len(formatArgs) > 0 {
		reason = fmt.Sprintf(reason, formatArgs)
	}

	return errors.NewErrorFromString("TxValidation FAILED: %s ", reason)
}

func TxValidationUnknownTransactionType(t string) error {
	return TxValidationFailed("unknown transaction type %s", t)
}
