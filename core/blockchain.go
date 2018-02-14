package core

type Blockchain struct {
	sink *BlockStore
}

func (self Blockchain) Cursor() (*BlockCursor, error) {
	var sink = (*self.sink)
	var result, err = sink.Cursor()
	return result, err
}

func (self *Blockchain) AddNewBlock(data string) error {
	var sink = (*self.sink)
	var last, err = sink.Tip()
	if err != nil {
		return err
	}
	if last == nil {
		last = NewGenesisBlock()
		sink.Add(last)
	}

	var newBlock = NewBlock(data, last.Hash)
	err = sink.Add(newBlock)

	return err
}

// exported
// type Blockchain interface {
// 	GetBlocks() []*Block
// 	AddNewBlock(data string)
// }

func NewBlockchain(sink *BlockStore) (*Blockchain, error) {

	var last, err = (*sink).Tip()
	if err != nil {
		return nil, err
	}

	if last == nil {
		last = NewGenesisBlock()
		(*sink).Add(last)
	}

	var self = new(Blockchain)
	self.sink = sink

	return self, nil
}
