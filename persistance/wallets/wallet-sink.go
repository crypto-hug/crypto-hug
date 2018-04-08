package wallets

import (

	// "../../formatters"
	// "../../serialization"

	"path/filepath"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/persistance/json-store"
	"github.com/spf13/afero"
)

type FsWalletSink struct {
	store *json_store.FsJsonStore
	root  string
}

func NewFsWalletSink(root string) *FsWalletSink {
	result := FsWalletSink{}

	result.store = json_store.NewFsJsonStore(root)
	result.root = root

	return &result
}

// func (self *FsWalletSink) PutMetadata(address *core.Address, key string, data string) error {
// 	path := filepath.Join(self.root, address.Address, "_meta.json")
// 	err := self.store.ModifyJsonFile(path, func(metadata *gabs.Container) (*gabs.Container, error) {
// 		if metadata == nil {
// 			metadata = gabs.New()
// 		}

// 		metadata, err := metadata.SetP(address, "address")
// 		if err != nil {
// 			return nil, errors.Wrap(err, "FsBlockSink:PutMetadata")
// 		}

// 		metadata, err = metadata.SetP(data, key)
// 		if err != nil {
// 			return nil, errors.Wrap(err, "FsBlockSink:PutMetadata")
// 		}

// 		return metadata, nil
// 	})

// 	return err
// }

// func (self *FsWalletSink) GetMetadata(address *core.Address, key string) (string, error) {
// 	path := filepath.Join(self.root, address.Address, "_meta.json")
// 	metadata, err := self.store.ReadJsonFile(path)
// 	if err != nil {
// 		return "", errors.Wrap(err, "FsBlockSink:GetMetadata")
// 	}
// 	if metadata == nil {
// 		return "", nil
// 	}

// 	result := metadata.Path(key).Data().(string)

// 	return result, nil
// }

func (self *FsWalletSink) assetPath(address *core.Address, asset *core.AssetHeader) string {
	path := filepath.Join(self.root, address.Address, string(asset.Type), asset.Address.Address+".json")
	return path
}

func (self *FsWalletSink) getAssetMetadata(address *core.Address, asset *core.AssetHeader, key string) (interface{}, error) {
	path := self.assetPath(address, asset)
	assetFile, err := self.store.ReadJsonFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read json file")
	}

	if assetFile == nil {
		return nil, nil
	}

	prop := assetFile.Path(key)
	if prop == nil {
		return nil, nil
	}
	result := prop.Data()

	return result, nil
}

func (self *FsWalletSink) putAssetMetadata(address *core.Address, asset *core.AssetHeader, key string, value interface{}) error {
	path := self.assetPath(address, asset)
	err := self.store.ModifyJsonFile(path, func(assetFile *gabs.Container) (*gabs.Container, error) {
		if assetFile == nil {
			assetFile = gabs.New()
		}

		_, err := assetFile.SetP(value, key)
		if err != nil {
			return nil, errors.Wrap(err, "could not set asset metadata")
		}

		return assetFile, nil
	})

	return err
}

func (self *FsWalletSink) PutAssetPropF(address *core.Address, asset *core.AssetHeader, key string, value float64) error {
	return self.putAssetMetadata(address, asset, key, value)
}

func (self *FsWalletSink) PutAssetPropS(address *core.Address, asset *core.AssetHeader, key string, value string) error {
	return self.putAssetMetadata(address, asset, key, value)
}

func (self *FsWalletSink) GetAssetPropF(address *core.Address, asset *core.AssetHeader, key string) (float64, error) {
	result, err := self.getAssetMetadata(address, asset, key)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, errors.PropertyNotExists
	}

	return result.(float64), nil
}

func (self *FsWalletSink) HasAsset(address *core.Address, asset *core.AssetHeader) (bool, error) {
	path := self.assetPath(address, asset)
	exists, err := afero.Exists(self.store.Fs(), path)

	return exists, err
}

func (self *FsWalletSink) RemoveAsset(address *core.Address, asset *core.AssetHeader) error {
	path := self.assetPath(address, asset)
	return self.store.Fs().Remove(path)
}

func (self *FsWalletSink) ListAssetsByType(address *core.Address, assetType core.AssetType) ([]string, error) {
	assetTypePath := filepath.Join(self.root, address.Address, string(assetType))
	exists, err := afero.DirExists(self.store.Fs(), assetTypePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not check asset type dir")
	}
	if exists == false {
		return []string{}, nil
	}

	files, err := afero.ReadDir(self.store.Fs(), assetTypePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not list asset types dir")
	}

	result := []string{}
	for _, file := range files {
		fileName := file.Name()
		assetAddr := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		result = append(result, assetAddr)
	}

	return result, nil
}

// func (self *FsWalletSink) PutAssetMetadataS(address *core.Address, asset *core.Asset, key string, value string) error {
// 	return self.putAssetMetadata(address, asset, key, value)
// }

// func (self *FsWalletSink) PutAssetMetadataN(address *core.Address, asset *core.Asset, key string, value string) error {
// 	return self.putAssetMetadata(address, asset, key, value)
// }

// func (self *FsWalletSink) PutBalance(address *core.Address, asset *core.Asset, newBalance int) error {
// 	path := filepath.Join(self.root, address.Address, string(asset.Type), asset.Address+".json")
// 	err := self.store.ModifyJsonFile(path, func(assetFile *gabs.Container) (*gabs.Container, error) {
// 		if assetFile == nil {
// 			assetFile = gabs.New()
// 		}

// 		_, err := assetFile.SetP(address, "address")
// 		if err != nil {
// 			return nil, errors.Wrap(err, "FsBlockSink:PutBalance")
// 		}

// 		_, err = assetFile.SetP(newBalance, "balance")
// 		if err != nil {
// 			return nil, errors.Wrap(err, "FsBlockSink:PutBalance")
// 		}

// 		return assetFile, errors.Wrap(err, "FsBlockSink:PutBalance")
// 	})

// 	return err

// }
