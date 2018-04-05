package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
)

type SpawnHugProcessor struct {
	sink core.WalletSink
}

func NewSpawnHugProcessor(sink core.WalletSink) *SpawnHugProcessor {
	result := SpawnHugProcessor{sink: sink}
	return &result
}

func (self SpawnHugProcessor) Validate(tx *core.Transaction) error {
	err := errors.MustBeNotNil("tx", tx)
	if err != nil {
		return err
	}

	if tx.Type != SpawnHugTxType {
		return errors.NewErrorFromString("Unknown tx type %s", tx.Type)
	}

	data, err := UnwrapSpawnHugTxData(tx)
	if err != nil {
		return errors.Wrap(err, "SpawnHugProcessor:Validate")
	}

	if data == nil {
		return errors.NewErrorFromString("tx content is not valid SpawnHugTxData")
	}

	return nil
}

func (self *SpawnHugProcessor) ShouldProcess(tx *core.Transaction) bool {
	return tx.Type == SpawnHugTxType
}

func (self *SpawnHugProcessor) Name() string {
	return "SpawnHug"
}

func (self *SpawnHugProcessor) Process(tx *core.Transaction) error {
	data, err := UnwrapSpawnHugTxData(tx)
	if err != nil {
		return errors.Wrap(err, "SpawnHugProcessor:Process")
	}

	balance, err := self.sink.GetBalance(data.RecipientAddress, &data.Asset)
	if err != nil {
		return errors.Wrap(err, "SpawnHugProcessor:Process")
	}

	if balance != 0 {
		return errors.NewErrorFromString("inavlid balance %v expected 0 for address %s", balance, data.RecipientAddress)
	}

	err = self.sink.PutBalance(data.RecipientAddress, &data.Asset, balance+1)

	return errors.Wrap(err, "SpawnHugProcessor:Process")
}
