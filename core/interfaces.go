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

type AssetStats interface {
	AssetType() string
	Count() int64
	Set(key string, data interface{}) error
	Get(key string, data interface{}) error
}

type WalletSink interface {
	PutMetadata(address string, key string, data string) error
	GetMetadata(address string, key string) (string, error)
	GetBalance(address string, asset *Asset) (int, error)
	PutBalance(address string, asset *Asset, newBalance int) error
}

type TransactionProcessor interface {
	ShouldProcess(tx *Transaction) bool
	Validate(tx *Transaction) (bool error)
	Process(tx *Transaction) error
	Name() string
}

type TransactionProcessors []TransactionProcessor

//type TxValidationIssue string

//type TxValidationIssues []TxValidationIssue
