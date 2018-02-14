package core

type BlockStore interface {
	Add(block *Block) error
	Tip() (*Block, error)
	//Prev(current *Block) (*Block, error)
	Cursor() (*BlockCursor, error)
}

type BlockCursor interface {
	Reset() error
	Next() (bool, error)
	Current() *Block
}
