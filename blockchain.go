package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/v-braun/go-must"
)

type Blockchain struct {
	config *Config
	fs     *fs.FileSystem
	proc   *TxProcessor
	states *StateStore
	store  *TxStore
}

func NewBlockchain(fs *fs.FileSystem, config *Config) *Blockchain {
	bc := new(Blockchain)
	bc.fs = fs
	bc.config = config
	bc.proc = NewTxProcessor(fs, config)
	bc.states = newStateStore(fs, config)
	bc.store = newTxStore(fs, config)

	return bc
}

func (bc *Blockchain) ProcessTransaction(tx *Transaction) error {
	err := bc.proc.Process(tx)
	return err
}

func (bc *Blockchain) CreateGenesisBlockIfNotExists() {
	genesisHugPubKey := utils.Base58FromStringMust(bc.config.GenesisTx.PubKey)
	genesisHugAddr, err := NewAddress(genesisHugPubKey)
	must.NoError(err, "failed create genesis hug address")

	if bc.states.HugExists(genesisHugAddr) {
		return
	}

	genesisTx := NewGenesisTransaction(bc.config)

	err = bc.proc.Process(genesisTx)
	must.NoError(err, "unexpected genesis block processing")

	_, err = bc.store.CommitStagedTx()
	must.NoError(err, "could not commit staged genesis tx")
}
