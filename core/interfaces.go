package core

// type BlockStore interface {
// 	Add(block *Block) error
// 	Tip() (*Block, error)
// 	GenesisBlock() (*Block, error)
// 	Cursor() (*BlockCursor, error)
// }

// type BlockCursor interface {
// 	Reset() error
// 	Next() (bool, error)
// 	Current() *Block
// }

type BlockSink interface {
	Put(block *Block) error
	Get(hash []byte) (*Block, error)
}
type BlockStats interface {
	PutTip(hash []byte) error
	PutGenesis(hash []byte) error
	GetTip() ([]byte, error)
	GetGenesis() ([]byte, error)
}

type AssetState interface {
	Get(addr []byte) (*Asset, error)
	Set(asset *Asset) error
}

type TxValidator interface {
	Validate(tx *Transaction) error
}

type TxValidators []TxValidator

type TxValidatorsRegistry interface {
	Get(tx *Transaction) (TxValidators, error)
}
