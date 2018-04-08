package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/serialization"
)

type SpendHugTxData struct {
	HugAddress       string
	RecipientAddress string
}

const (
	SpawnHugTxType core.TransactionType = "SPAWNHUG"
	SpendHugTxType core.TransactionType = "SPENDHUG"
)

var DataEmptyError error = errors.NewErrorFromString("No Data")

func NewSpawnHugTransaction(producer *core.Address, pubKey []byte) (*core.Transaction, error) {
	var result = core.NewTransaction(SpawnHugTxType, producer.Address, pubKey, []byte{})
	return result, nil
}

func NewSpendHugTransaction(senderAddress string, senderPubKey []byte, hugAddress string, recipientAddress string) (*core.Transaction, error) {
	data := SpendHugTxData{HugAddress: hugAddress, RecipientAddress: recipientAddress}
	dataRaw, err := serialization.ObjToJsonRaw(data)
	if err != nil {
		panic(err)
	}

	tx := core.NewTransaction(SpendHugTxType, senderAddress, senderPubKey, dataRaw)
	return tx, nil
}

func unwrap(tx *core.Transaction, content interface{}) error {
	if tx.Data == nil {
		return DataEmptyError
	}
	if len(tx.Data) <= 0 {
		return DataEmptyError
	}

	err := serialization.JsonParseRaw(tx.Data, content)
	return err
}

func UnwrapSpendHugTxData(tx *core.Transaction) (*SpendHugTxData, error) {
	result := &SpendHugTxData{}
	err := unwrap(tx, result)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse pend hug data")
	}
	if result.RecipientAddress == "" {
		return nil, errors.NewErrorFromString("no recipient address defined")
	}
	if result.HugAddress == "" {
		return nil, errors.NewErrorFromString("no hug address defined")
	}

	if _, err := core.NewAddressFromString(result.RecipientAddress); err != nil {
		return nil, errors.Wrap(err, "could not parse receipient address string")
	}

	if _, err = core.NewAddressFromString(result.HugAddress); err != nil {
		return nil, errors.Wrap(err, "could not parse hug address string")
	}

	return result, err
}
