package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
)

type TxStore struct {
	fs   *fs.FileSystem
	conf *Config
}

func newTxStore(fs *fs.FileSystem, conf *Config) *TxStore {
	result := new(TxStore)
	result.fs = fs
	result.conf = conf
	return result
}
func (s *TxStore) readTx(path string) (*Transaction, error) {
	content, err := s.fs.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tx := new(Transaction)
	if err := utils.JsonParseRaw(content, tx); err != nil {
		return nil, errors.Wrapf(err, "could not json parse tx in file %s content: %s", path, string(content))
	}

	return tx, nil
}

func (s *TxStore) commitBlock(block *Block) error {
	addr, err := NewAddress(block.Hash)
	if err != nil {
		return errors.Wrap(err, "could not create address for block")
	}

	path := s.conf.Paths.BlockDir + addr + ".json"
	if s.fs.FileExists(path) {
		return errors.Errorf("invalid block commit [%s]. Block already exists")
	}

	blockData, err := utils.JsonSerializeRaw(block)
	if err != nil {
		return errors.Wrapf(err, "could not serialize block [%s]", addr)
	}

	err = s.fs.WriteFile(path, blockData)
	if err != nil {
		return errors.Wrapf(err, "could not write block [%s]", addr)
	}

	for i, tx := range block.Transactions {
		addr, err := NewAddress(tx.Hash.Bytes())
		if err != nil {
			return errors.Wrapf(err, "could not create address for tx (%d) in block %s", i, addr)
		}

		txPath := s.conf.Paths.TxStagePath + addr + ".json"
		err = s.fs.Remove(txPath)
		if err != nil {
			return errors.Wrapf(err, "could not remove staged tx (%s)", txPath)
		}
	}

	return nil
}

func (s *TxStore) CommitStagedTx() (*Block, error) {
	stagePath := s.conf.Paths.TxStagePath
	files, err := s.fs.ListDir(stagePath)
	if err != nil {
		return nil, err
	}

	if len(files) == 1 && s.fs.FileNameWithoutExt(files[0].Name()) == s.conf.GenesisTx.Address {
		genPath := stagePath + files[0].Name()
		tx, err := s.readTx(genPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed read tx file %s", genPath)
		}

		genesisBlock := NewGenesisBlock(s.conf, tx)
		err = s.commitBlock(genesisBlock)

		return genesisBlock, err
	}

	block := NewBlock()
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filePath := stagePath + f.Name()

		content, err := s.fs.ReadFile(filePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed read tx file %s", filePath)
		}

		tx := new(Transaction)
		if err := utils.JsonParseRaw(content, tx); err != nil {
			return nil, errors.Wrapf(err, "could not json parse tx in file %s", filePath)
		}

		block.Transactions = append(block.Transactions, tx)
	}

	return block, nil
}

func (s *TxStore) StageTx(ctx *txProcessCtx) error {
	path := s.conf.Paths.TxStagePath + ctx.address + ".json"
	if s.fs.FileExists(path) {
		return errors.Errorf("tx [%s] already staged", ctx.address)
	}

	data, err := utils.JsonSerializeRaw(ctx.tx)
	if err != nil {
		return errors.Wrapf(err, "could not serialize tx [%s]", ctx.address)
	}

	err = s.fs.WriteFile(path, data)
	return errors.Wrapf(err, "could not stage tx to file [%s]", path)
}
