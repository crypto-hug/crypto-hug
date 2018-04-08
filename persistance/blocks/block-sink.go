package blocks

import (
	"os"
	"path/filepath"

	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/formatters"
	"github.com/crypto-hug/crypto-hug/serialization"
	"github.com/spf13/afero"
)

type FsBlockSink struct {
	fs   afero.Fs
	root string
}

func NewFsBlockSink(root string) *FsBlockSink {
	result := FsBlockSink{}

	if len(root) == 0 {
		result.fs = afero.NewMemMapFs()
	} else {
		result.fs = afero.NewOsFs()

	}

	result.root = root
	// todo: check why 0600 or 0666 cause permission denied for CreateFile
	result.fs.MkdirAll(result.root, 0700)

	return &result
}

func (self *FsBlockSink) pathToBlockHash(hash []byte) (string, error) {
	file := formatters.HexStringFromRaw(hash)
	file = filepath.Join(self.root, file)
	file, err := filepath.Abs(file)
	return file, err
}

func (self *FsBlockSink) Put(block *core.Block) error {
	file, err := self.pathToBlockHash(block.Hash)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:Put:GetFilePath")
	}

	bin, err := serialization.ObjEncode(block)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:Put:EncodeBlock")
	}

	handle, err := self.fs.Create(file)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:Put:CreateFile")
	}
	defer handle.Close()

	_, err = handle.Write(bin)

	return err
}

func (self *FsBlockSink) Get(hash []byte) (*core.Block, error) {
	file, err := self.pathToBlockHash(hash)
	if err != nil {
		return nil, errors.Wrap(err, "FsBlockSink:Get:GetFilePath")
	}

	handle, err := self.fs.Open(file)
	if err == os.ErrNotExist {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "FsBlockSink:Get:OpenFile")
	}
	defer handle.Close()

	bin, err := afero.ReadAll(handle)
	if err != nil {
		return nil, errors.Wrap(err, "FsBlockSink:Get:ReadFile")
	}

	result := new(core.Block)
	err = serialization.ObjDecode(bin, result)

	return result, errors.Wrap(err, "FsBlockSink:Get:ObjDecode")
}
