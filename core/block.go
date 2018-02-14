package core

import (
	"fmt"
	"time"

	"../common/crypt"
	"../common/formatters"
)

// unexported

type Block struct {
	Version   uint16
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

func (self *Block) calcHash() []byte {

	var hash = crypt.AllBytesHash(
		[]byte(formatters.HexString(int64(self.Version))),
		self.PrevHash,
		[]byte(formatters.HexString(self.Timestamp)),
		self.Data)

	return hash
}

func (self *Block) PrettyPrint() string {
	const tmpl = `
Version:	%d
Hash:		%x
Prev. Hash:	%x
Timestamp:	%d
Data:		%s
`
	var result = fmt.Sprintf(tmpl,
		self.Version, self.Hash, self.PrevHash,
		self.Timestamp, self.Data)

	return result
}

// exported
const CURRENT_VERSION uint16 = 1

// type Block interface {
// 	GetVersion() uint16
// 	GetData() string
// 	GetHash() []byte
// }

func NewBlock(data string, prevHash []byte) *Block {
	var self = new(Block)
	self.Version = CURRENT_VERSION
	self.Timestamp = time.Now().Unix()
	self.Data = []byte(data)
	self.PrevHash = prevHash
	self.Hash = self.calcHash()

	return self
}
func NewGenesisBlock() *Block {
	var self = NewBlock("Genesis Block", []byte{})
	return self
}
