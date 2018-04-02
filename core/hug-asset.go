package core

import (
	"../serialization"
	"./errors"
)

type HugHistoryEntry struct {
	Action         HugAction
	OriginatorAddr []byte
	EntryDate      int64
}

type HugHistory []HugHistoryEntry

type HugState string

const (
	HugStateNew   HugState = "NEW"
	HugStateDeath HugState = "DEATH"
	HugStateOwned HugState = "OWNED"
)

type HugAction string

const (
	HugActionBirth       HugAction = "BRTH"
	HugActionReincarnate HugAction = "RENCRNTE"
	HugActionDonate      HugAction = "DNTE"
	HugActionDie         HugAction = "DIE"
)

type Hug struct {
	History HugHistory
}

func UnwrapHug(asset *Asset) (*Hug, error) {
	if asset.Type != AssetTypeHug {
		return nil, errors.NewErrorFromString("Invalid Asset type %s. Expected %s | Asset: %s", asset.Type, AssetTypeHug, asset.Address)
	}

	if len(asset.Data) <= 0 {
		return nil, errors.NewErrorFromString("Invalid Data content on Asset %s", asset.Address)
	}

	hug := Hug{}
	err := serialization.JsonParse(asset.Data, hug)

	return &hug, err
}
