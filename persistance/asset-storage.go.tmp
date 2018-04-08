package hugdb

import (
	"../core"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
)

type BoltAssetState struct {
	db *BoltDb
}

func NewBoltAssetState(db *BoltDb) *BoltAssetState {
	var result = BoltAssetState{db: db}
	return &result
}

func (self *BoltAssetState) Get(addr []byte) (*core.Asset, error) {
	var result *core.Asset = nil

	err := self.db.GetDb().View(func(tx *bolt.Tx) error {
		var existing = new(core.Asset)
		var has, err = self.db.GetDecode(tx, self.db.BucketAssets(), addr, &existing)
		if has {
			result = existing
		}
		return err
	})

	return result, err
}

func (self *BoltAssetState) Set(asset *core.Asset) error {
	err := self.db.GetDb().Update(func(tx *bolt.Tx) error {
		if asset.Etag != core.NewEtag {
			var existing = new(core.Asset)
			var has, err = self.db.GetDecode(tx, self.db.BucketAssets(), asset.AddressRaw(), &existing)
			if err != nil {
				return err
			}
			if has == false {
				return errors.New(fmt.Sprintf("could not find asset %s with etag %v", asset.Address, asset.Etag))
			}
			if existing.Etag != asset.Etag {
				return errors.New(fmt.Sprintf("etag missmatch for asset %s. expected etag: %v | current etag: %v", asset.Address, asset.Etag, existing.Etag))
			}
			asset.Etag = asset.Etag + 1
			err = self.db.SetEncoded(tx, self.db.BucketAssets(), asset.AddressRaw(), asset)
			return err
		} else {
			asset.Etag = asset.Etag + 1
			err := self.db.SetEncoded(tx, self.db.BucketAssets(), asset.AddressRaw(), asset)
			return err
		}
	})

	return err
}
