package chug

import (
	"../"
	"../../serialization"
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
