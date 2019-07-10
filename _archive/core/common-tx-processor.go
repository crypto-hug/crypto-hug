package core

import (
	"bytes"

	"github.com/crypto-hug/crypto-hug/errors"
)

type CommonTxProcessor struct {
}

func (self CommonTxProcessor) Setup(wallets WalletSink, assets AssetSink) {
}

func (self CommonTxProcessor) Prepare(tx *Transaction) error {

	if err := errors.MustBeNotNil("tx", tx); err != nil {
		return err
	}

	if bytes.Compare(tx.Hash, tx.CalcHash()) != 0 {
		return TxValidationFailed("Tx hash missmatch")
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
