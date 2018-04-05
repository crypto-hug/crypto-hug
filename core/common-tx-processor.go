package core

import (
	"bytes"

	"github.com/crypto-hug/crypto-hug/errors"
)

type CommonTxProcessor struct {
}

func (self *CommonTxProcessor) Validate(tx *Transaction) error {
	err := errors.MustBeNotNil("tx", tx)
	if err != nil {
		return err
	}

	if bytes.Compare(tx.Hash, tx.CalcHash()) != 0 {
		return TxValidationFailed("Tx hash missmatch")
	}

	if len(tx.Data) == 0 {
		return TxValidationFailed("tx data is empty")
	}
	return nil
}

func (self *CommonTxProcessor) Name() string {
	return "Common"
}

func (self *CommonTxProcessor) ShouldProcess(tx *Transaction) bool {
	return true
}

func (self *CommonTxProcessor) Process(tx *Transaction) error {
	return nil
}
