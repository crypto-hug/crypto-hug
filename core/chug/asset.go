package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
)

const AssetTypeHug core.AssetType = "HUG"

// type HugHistoryEntry struct {
// 	Action         HugAction
// 	OriginatorAddr []byte
// 	EntryDate      int64
// }

// type HugHistory []HugHistoryEntry

// type HugState string

// const (
// 	HugStateNew   HugState = "NEW"
// 	HugStateDeath HugState = "DEATH"
// 	HugStateOwned HugState = "OWNED"
// )

const (
	HugActionBirth core.JournalAction = "BIRTH"

	// 	HugActionReincarnate HugAction = "RENCRNTE"
	HugActionSpend core.JournalAction = "SPEND"

// 	HugActionDie         HugAction = "DIE"
)

type HugAsset struct {
	header   *core.AssetHeader
	producer *core.Address
}

func NewHugAsset(producer *core.Address) *HugAsset {
	result := HugAsset{}
	result.producer = producer
	result.header = core.NewAssetHeader(result.producer, AssetTypeHug)
	return &result
}
func NewExistingHugAddress(producer *core.Address, address, assetAddress *core.Address) {
	result := NewHugAsset(producer)
	result.header.Address = assetAddress
}

func (self *HugAsset) Header() *core.AssetHeader {
	return self.header
}

func (self *HugAsset) Bytes() []byte {
	return []byte{}
}
