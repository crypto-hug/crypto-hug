package core

import (
	"../formatters"
	"./errors"
)

type AssetType string

const (
	AssetTypeHug AssetType = "HUG"
	NewEtag      int32     = 0
)

type Asset struct {
	Version         Version
	Address         string
	ProducerAddress string
	Etag            int32
	Type            AssetType
	Data            string
}

func NewHugAsset(producer *Address) (*Asset, error) {
	addr, err := NewAddress()
	if err != nil {
		return nil, err
	}

	result := Asset{Version: AssetVersion,
		Address:         addr.Address,
		ProducerAddress: producer.Address,
		Etag:            NewEtag,
		Type:            AssetTypeHug,
		Data:            ""}

	return &result, nil
}

func (self *Asset) AddressRaw() []byte {
	raw, err := formatters.Base58FromString(self.Address)
	if err != nil {
		panic(errors.Wrap(err, "Asset:AddressRaw"))
	}

	return raw
}
