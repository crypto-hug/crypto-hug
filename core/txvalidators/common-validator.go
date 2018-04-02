package txvalidators

import (
	".."
	"../errors"
	"bytes"
)

type CommonValidator struct {
}

func (self *CommonValidator) Validate(tx *core.Transaction) error {
	err := errors.MustBeNotNil("tx", tx)
	if err != nil {
		return err
	}

	if bytes.Compare(tx.Hash, tx.CalcHash()) != 0 {
		return errors.TxValidationFailed("Tx hash missmatch")
	}

	if len(tx.Data) == 0 {
		return errors.TxValidationFailed("tx data is empty")
	}
	return nil
}
