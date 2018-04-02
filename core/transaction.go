package core

import (
	"../crypt"
	"../formatters"
	"../serialization"
	"fmt"
	"time"
)

type TransactionType string

const (
	SpawnHugTxType  TransactionType = "SPAWNHUG"
	DonateHugTxType TransactionType = "DONATEHUG"
)

type Transactions []*Transaction

type Transaction struct {
	Version   Version
	Hash      []byte
	Type      TransactionType
	Timestamp int64
	Data      []byte
}

type SpawnHugTxData struct {
	RecipientAddress string
	Asset            Asset
}

func (self Transactions) getHash() []byte {
	var all [][]byte
	for _, tx := range self {
		var hash = tx.Hash
		all = append(all, hash)
	}

	var result = crypt.AllBytesHash(all[:]...)
	return result
}

func (self *Transaction) CalcHash() []byte {
	var result = crypt.AllBytesHash(
		[]byte(formatters.HexString(int64(self.Version))),
		[]byte(self.Type),
		[]byte(formatters.HexString(self.Timestamp)),
		self.Data,
	)

	return result
}

func (self *Transaction) PrettyPrint() string {
	const tmpl = `
Version:        %d
Hash:           %x
Timestamp:      %d
Type:           %v
Data:           %v`
	var result = fmt.Sprintf(tmpl,
		self.Version,
		self.Hash,
		self.Timestamp,
		self.Type,
		string(self.Data))

	return result
}

func NewTransaction(ofType TransactionType, withData []byte) *Transaction {
	var result = Transaction{Version: TxVersion,
		Type:      ofType,
		Data:      withData,
		Timestamp: time.Now().Unix()}

	result.Hash = result.CalcHash()
	return &result
}

func NewSpawnHugTransaction(producerAdr *Address) (*Transaction, error) {
	asset, err := NewHugAsset(producerAdr)
	if err != nil {
		return nil, err
	}

	var txData = SpawnHugTxData{RecipientAddress: asset.ProducerAddress, Asset: *asset}
	txDataRaw, err := serialization.ObjToJsonRaw(txData)
	if err != nil {
		return nil, err
	}

	var result = NewTransaction(SpawnHugTxType, txDataRaw)
	return result, nil
}
