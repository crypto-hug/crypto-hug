package core

import (
	"fmt"
	"time"

	"github.com/crypto-hug/crypto-hug/crypt"
	"github.com/crypto-hug/crypto-hug/formatters"
)

type TransactionType string

type Transactions []*Transaction

type Transaction struct {
	Version   Version
	Hash      []byte
	Type      TransactionType
	Timestamp int64
	Data      []byte
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
