package wallets

import (
	// "../../formatters"
	// "../../serialization"
	"io/ioutil"
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/spf13/afero"
)

var metadataTemplate = []byte(`{
	  "address": ""
}`)
var assetFileTemplate = []byte(`{
	"address": "",
	"balance": 0
}`)

type FsWalletSink struct {
	fs   afero.Fs
	root string
}

func NewFsWalletSink(root string) *FsWalletSink {
	result := FsWalletSink{}

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
func (self *FsWalletSink) readJsonFile(path string) ([]byte, error) {
	self.fs.MkdirAll(filepath.Dir(path), 0700)

	handle, err := self.fs.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "FsBlockSink:readJsonFile")
	}
	defer handle.Close()

	bin, err := ioutil.ReadAll(handle)
	if err != nil {
		return nil, errors.Wrap(err, "FsBlockSink:readJsonFile")
	}
	if len(bin) <= 0 {
		return nil, nil
	}
	return bin, nil
}

func (self *FsWalletSink) writeJsonFile(path string, file []byte) error {
	self.fs.MkdirAll(filepath.Dir(path), 0700)

	handle, err := self.fs.Create(path)

	if err != nil {
		return errors.Wrap(err, "FsBlockSink:writeJsonFile")
	}
	defer handle.Close()

	_, err = handle.Write(file)
	return err
}

func (self *FsWalletSink) writeMetadataFile(address string, file []byte) error {
	metafile := filepath.Join(self.root, address, "_meta.json")
	return self.writeJsonFile(metafile, file)
}
func (self *FsWalletSink) readMetadataFile(address string) ([]byte, error) {
	metafile := filepath.Join(self.root, address, "_meta.json")
	return self.readJsonFile(metafile)
}

func (self *FsWalletSink) writeAssetFile(address string, asset *core.Asset, file []byte) error {
	assetFile := filepath.Join(self.root, address, string(asset.Type), asset.Address+".json")
	return self.writeJsonFile(assetFile, file)
}
func (self *FsWalletSink) readAssetFile(address string, asset *core.Asset) ([]byte, error) {
	assetFile := filepath.Join(self.root, address, string(asset.Type), asset.Address+".json")
	return self.readJsonFile(assetFile)
}

func (self *FsWalletSink) PutMetadata(address string, key string, data string) error {
	metadata, err := self.readMetadataFile(address)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutMetadata")
	}
	if metadata == nil {
		metadata = metadataTemplate
	}

	metadataJson, err := gabs.ParseJSON(metadata)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutMetadata")
	}

	metadataJson, err = metadataJson.SetP(address, "address")
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutMetadata")
	}

	metadataJson, err = metadataJson.SetP(data, key)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutMetadata")
	}

	err = self.writeMetadataFile(address, metadataJson.Bytes())

	return errors.Wrap(err, "FsBlockSink:PutMetadata")
}

func (self *FsWalletSink) GetMetadata(address string, key string) (string, error) {
	metadata, err := self.readMetadataFile(address)
	if err != nil {
		return "", errors.Wrap(err, "FsBlockSink:GetMetadata")
	}
	if metadata == nil {
		return "", nil
	}

	metadataJson, err := gabs.ParseJSON(metadata)
	if err != nil {
		return "", errors.Wrap(err, "FsBlockSink:GetMetadata")
	}

	result := metadataJson.Path(key).Data().(string)

	return result, nil
}

func (self *FsWalletSink) GetBalance(address string, asset *core.Asset) (int, error) {
	assetFile, err := self.readAssetFile(address, asset)
	if err != nil {
		return 0, errors.Wrap(err, "FsBlockSink:GetBalance")
	}

	if assetFile == nil {
		return 0, nil
	}

	assetFileJson, err := gabs.ParseJSON(assetFile)
	if err != nil {
		return 0, errors.Wrap(err, "FsBlockSink:GetBalance")
	}

	val := assetFileJson.Path("balance").Data().(int)

	return val, nil
}
func (self *FsWalletSink) PutBalance(address string, asset *core.Asset, newBalance int) error {
	assetFile, err := self.readAssetFile(address, asset)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutBalance")
	}

	if assetFile == nil {
		assetFile = assetFileTemplate
	}

	assetFileJson, err := gabs.ParseJSON(assetFile)
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutBalance")
	}

	_, err = assetFileJson.SetP(address, "address")
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutBalance")
	}

	_, err = assetFileJson.SetP(newBalance, "balance")
	if err != nil {
		return errors.Wrap(err, "FsBlockSink:PutBalance")
	}

	err = self.writeAssetFile(address, asset, assetFileJson.Bytes())

	return errors.Wrap(err, "FsBlockSink:PutBalance")
}
