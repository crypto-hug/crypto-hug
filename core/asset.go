package core

import (
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/formatters"
)

type AssetType string

const (
	NewEtag int32 = 0
)

type Asset struct {
	Version         Version
	Address         string
	ProducerAddress string
	Etag            int32
	Type            AssetType
	Data            string
}

func (self *Asset) AddressRaw() []byte {
	raw, err := formatters.Base58FromString(self.Address)
	if err != nil {
		panic(errors.Wrap(err, "Asset:AddressRaw"))
	}

	return raw
}
