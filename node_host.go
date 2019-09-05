package chug

import (
	"sync"

	"github.com/crypto-hug/crypto-hug/fs"
)

type NodeHost struct {
	bc   *Blockchain
	lock sync.Locker
}

func NewNodeHost(fs *fs.FileSystem, config *Config) *NodeHost {
	result := new(NodeHost)
	result.bc = NewBlockchain(fs, config)
	result.lock = new(sync.Mutex)

	return result
}

func (nh *NodeHost) Start() {
	nh.lock.Lock()
	defer nh.lock.Unlock()

	nh.bc.CreateGenesisBlockIfNotExists()
	// TODO:
	// setup network
	//		1. connect to enough peers
	// sync bc
	//		1. get longest blockchain
	//		2. download missing blocks & process them
	//		3. staged tx are block id "nil"
	// begin listening for IO
	//		1. listen for transactions
	//		2. listen for discover requests
	// frequentley check bc in sync
	// 		1. check current latest block with most latest in network (min x% nodes)
	//		2. pause processing
	// 		3. begin sync again all block headers

	// assume everything is now started
}

func (nh *NodeHost) ProcessTransaction(tx *Transaction) error {
	nh.lock.Lock()
	defer nh.lock.Unlock()

	err := nh.bc.ProcessTransaction(tx)
	return err
}
