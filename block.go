package chug

import (
	"bytes"
	"time"

	"github.com/crypto-hug/crypto-hug/utils"
)

type Block struct {
	Version      Version
	Hash         []byte
	PrevHash     []byte
	Timestamp    int64
	Transactions Transactions
}

func NewBlock() *Block {
	b := new(Block)
	b.Version = BlockVersion
	b.Timestamp = time.Now().Unix()
	return b
}

func NewGenesisBlock(config *Config, genesisTx *Transaction) *Block {
	b := NewBlock()
	b.Timestamp = genesisTx.Timestamp
	b.Transactions = append(b.Transactions, genesisTx)
	b.HashBlock()
	return b
}

func (b *Block) getHash() []byte {
	data := bytes.Join([][]byte{
		[]byte(b.Version),
		[]byte(b.PrevHash),
		[]byte(utils.Int64GetBytes(b.Timestamp)),
		b.Transactions.getHash(),
	}, []byte{})

	result := utils.Hash(data)
	return result
}

func (b *Block) HashBlock() {
	b.Hash = b.getHash()
}
