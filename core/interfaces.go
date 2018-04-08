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

// type AssetStats interface {
// 	AssetType() string
// 	Count() int64
// 	Set(key string, data interface{}) error
// 	Get(key string, data interface{}) error
// }

type WalletSink interface {

	// PutS(address *Address, key string, value string) error
	// GetS(address *Address, key string) (float64, error)

	// PutF(address *Address, key string, value float64) error
	// GetF(address *Address, key string) (float64, error)

	HasAsset(address *Address, asset *AssetHeader) (bool, error)
	ListAssetsByType(address *Address, assetType AssetType) ([]string, error)
	RemoveAsset(address *Address, asset *AssetHeader) error
	PutAssetPropF(address *Address, asset *AssetHeader, key string, value float64) error
	PutAssetPropS(address *Address, asset *AssetHeader, key string, value string) error
	GetAssetPropF(address *Address, asset *AssetHeader, key string) (float64, error)
}

// type WalletStoredAsset interface {
// 	PutS(key string, value string) error
// 	PutN(key string, value float64) error
// 	GetS(key string) (string, error)
// 	GetN(key string) (float64, error)
// 	AppendS(key string, value string) error
// 	AppendN(key string, value float64) error
// }

type AssetJournalData map[string]interface{}

type JournalAction string

type AssetSink interface {
	PutJournal(asset *AssetHeader, action JournalAction, producer *Address, data *AssetJournalData) error
	GetHeader(address *Address) (*AssetHeader, error)
}

type TransactionProcessor interface {
	ShouldProcess(tx *Transaction) bool
	Setup(wallets WalletSink, assets AssetSink)
	Prepare(tx *Transaction) error
	Process(tx *Transaction) error
	Name() string
}

type TransactionProcessors []TransactionProcessor

//type TxValidationIssue string

//type TxValidationIssues []TxValidationIssue
