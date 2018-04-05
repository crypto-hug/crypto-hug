package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/serialization"
)

type SpawnHugTxData struct {
	RecipientAddress string
	Asset            core.Asset
}

const (
	SpawnHugTxType  core.TransactionType = "SPAWNHUG"
	DonateHugTxType core.TransactionType = "DONATEHUG"
)

func NewSpawnHugTransaction(producerAdr *core.Address) (*core.Transaction, error) {
	asset, err := NewHugAsset(producerAdr)
	if err != nil {
		return nil, err
	}

	var txData = SpawnHugTxData{RecipientAddress: asset.ProducerAddress, Asset: *asset}
	txDataRaw, err := serialization.ObjToJsonRaw(txData)
	if err != nil {
		return nil, err
	}

	var result = core.NewTransaction(SpawnHugTxType, txDataRaw)
	return result, nil
}

func UnwrapSpawnHugTxData(tx *core.Transaction) (*SpawnHugTxData, error) {
	if tx.Data == nil {
		return nil, nil
	}
	if len(tx.Data) <= 0 {
		return nil, nil
	}

	result := &SpawnHugTxData{}
	err := serialization.JsonParseRaw(tx.Data, result)

	return result, err
}
