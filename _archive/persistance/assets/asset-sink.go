package assets

import (
	"path/filepath"

	"github.com/Jeffail/gabs"
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/persistance/json-store"
)

type FsAssetSink struct {
	store *json_store.FsJsonStore
	root  string
}

func NewFsAssetSink(root string) *FsAssetSink {
	result := FsAssetSink{root: root}

	result.store = json_store.NewFsJsonStore(root)
	result.root = root

	return &result
}

func modifyStrict(container *gabs.Container, err error) {
	if err != nil {
		panic(err)
	}
}

func initAssetFile(file *gabs.Container, asset *core.AssetHeader) {
	modifyStrict(file.Object("header"))
	modifyStrict(file.SetP(asset.Version, "header.version"))
	modifyStrict(file.SetP(asset.Address.Address, "header.address"))
	modifyStrict(file.SetP(asset.Producer.Address, "header.producer"))
	modifyStrict(file.SetP(asset.Type, "header.type"))
	modifyStrict(file.Array("journal"))
}

func (self *FsAssetSink) path(address *core.Address) string {
	path := filepath.Join(self.root, address.Address+".json")
	return path
}

func (self *FsAssetSink) GetHeader(address *core.Address) (*core.AssetHeader, error) {
	path := self.path(address)
	file, err := self.store.ReadJsonFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not read asset file")
	}
	if file == nil {
		return nil, core.AssetNotExist
	}

	version := file.Path("header.version").Data().(float64)
	assetAddr := file.Path("header.address").Data().(string)
	producerAddr := file.Path("header.producer").Data().(string)
	assetType := file.Path("header.type").Data().(string)

	result := core.AssetHeader{}
	result.Address = core.NewAddressFromStringStrict(assetAddr)
	result.Version = core.Version(version)
	result.Producer = core.NewAddressFromStringStrict(producerAddr)
	result.Type = core.AssetType(assetType)

	return &result, nil
}

func (self *FsAssetSink) PutJournal(asset *core.AssetHeader, action core.JournalAction, producer *core.Address, data *core.AssetJournalData) error {
	path := self.path(asset.Address)
	err := self.store.ModifyJsonFile(path, func(file *gabs.Container) (*gabs.Container, error) {
		if file == nil {
			file = gabs.New()
			initAssetFile(file, asset)
		}

		content := map[string]interface{}{
			"action":   action,
			"producer": producer.Address,
			"data":     data}

		err := file.ArrayAppend(content, "journal")
		if err != nil {
			return nil, errors.Wrap(err, "could not write journal data")
		}

		return file, nil
	})

	return errors.Wrap(err, "could not update asset journal")
}
