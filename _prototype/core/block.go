package core

import (
	"fmt"
	"time"

	"../common/crypt"
	"../common/formatters"
)

// unexported

type Block struct {
	Version      uint16
	Timestamp    int64
	Hash         []byte
	PrevHash     []byte
	Transactions Transactions
	Nonce        int
}

func (self *Block) PrettyPrint(algo ProofAlgorithm) string {
	const tmpl = `
Version:        %d
Hash:           %x
Prev. Hash:     %x
Nonce:          %v
Timestamp:      %d
Valid:          %v
Transactions:   %v
`
	var result = fmt.Sprintf(tmpl,
		self.Version, self.Hash, self.PrevHash,
		self.Nonce,
		self.Timestamp,
		algo.Validate(self),
		len(self.Transactions))

	return result
}

// exported
const CurrentVersion uint16 = 1

// type Block interface {
// 	GetVersion() uint16
// 	GetData() string
// 	GetHash() []byte
// }

func NewBlock(transactions []*Transaction, prevHash []byte) *Block {
	var self = new(Block)
	self.Version = CurrentVersion
	self.Timestamp = time.Now().Unix()
	self.Transactions = transactions
	self.PrevHash = prevHash

	return self
}
func NewGenesisBlock(coinbase *Transaction) *Block {
	var self = NewBlock([]*Transaction{coinbase}, []byte{})
	self.Hash = []byte{}
	return self
}

func (self *Block) IsGenesisBlock() bool {
	return len(self.PrevHash) <= 0
}

func (self *Block) getTransactionsHash() {

}

func (self *Block) GetHash(proofTargetBits int) []byte {
	var result = crypt.AllBytesHash(
		self.PrevHash,
		self.Transactions.getHash(),
		[]byte(formatters.HexString(self.Timestamp)),
		[]byte(formatters.HexString(int64(proofTargetBits))),
		[]byte(formatters.HexString(int64(self.Nonce))),
	)

	return result
}
