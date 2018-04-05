package core

import (
	"fmt"
	"time"

	"github.com/crypto-hug/crypto-hug/crypt"
	"github.com/crypto-hug/crypto-hug/formatters"
)

type Block struct {
	Version      Version
	Hash         []byte
	PrevHash     []byte
	Timestamp    int64
	Transactions Transactions
}

func (self *Block) calcHash() []byte {
	// todo: calc a merkle root hash for the transactions
	var result = crypt.AllBytesHash(
		[]byte(formatters.HexString(self.Timestamp)),
		self.PrevHash,
		[]byte(formatters.HexString(self.Timestamp)),
		self.Transactions.getHash(),
	)

	return result
}

func NewBlock(transactions Transactions, prevHash []byte) *Block {
	var self = new(Block)
	self.Version = BlockVersion
	self.Timestamp = time.Now().Unix()
	self.Transactions = transactions
	self.PrevHash = prevHash
	self.Hash = self.calcHash()

	return self
}

func NewGenesisBlock(transactions Transactions) *Block {
	return NewBlock(transactions, []byte{})
}

func (self *Block) IsGenesisBlock() bool {
	return len(self.PrevHash) <= 0
}

func (self *Block) PrettyPrint() string {
	const tmpl = `
Version:        %d
Hash:           %x
Prev. Hash:     %x
Timestamp:      %d`
	var result = fmt.Sprintf(tmpl,
		self.Version, self.Hash, self.PrevHash,
		self.Timestamp)

	return result
}
