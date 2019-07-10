package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
)

const max_hugs_in_wallet = 5

type SpawnHugProcessor struct {
	wallets core.WalletSink
	assets  core.AssetSink

	receipient *core.Address
	asset      *HugAsset
}

func (self *SpawnHugProcessor) Setup(wallets core.WalletSink, assets core.AssetSink) {
	errors.AssertNotNil("wallets", wallets)
	errors.AssertNotNil("assets", assets)

	self.wallets = wallets
	self.assets = assets
}

func (self *SpawnHugProcessor) Prepare(tx *core.Transaction) error {
	if err := errors.MustBeNotNil("tx", tx); err != nil {
		return err
	}
	errors.AssertNotNil("self.wallets", self.wallets)
	errors.AssertNotNil("self.assets", self.assets)
	errors.AssertTrue("tx.Type", tx.Type == SpawnHugTxType, "is not SpawnHug")

	addr, err := core.NewAddressFromString(tx.Sender)
	if err != nil {
		return errors.Wrap(err, "invalid sender address")
	}

	asset := NewHugAsset(addr)

	ownedAssets, err := self.wallets.ListAssetsByType(addr, AssetTypeHug)
	if err != nil {
		return errors.Wrap(err, "could not get owned assets")
	}

	if ownedAssets != nil && len(ownedAssets) >= max_hugs_in_wallet {
		return errors.NewErrorFromString("max hug balance of %v reached: current %v", max_hugs_in_wallet, len(ownedAssets))
	}

	self.receipient = addr
	self.asset = asset

	return nil
}

func (self *SpawnHugProcessor) ShouldProcess(tx *core.Transaction) bool {
	return tx.Type == SpawnHugTxType
}

func (self *SpawnHugProcessor) Name() string {
	return "SpawnHug"
}

func (self *SpawnHugProcessor) Process(tx *core.Transaction) error {
	if err := self.wallets.PutAssetPropF(self.receipient, self.asset.Header(), "balance", 1); err != nil {
		return errors.Wrap(err, "failed to set new balance")
	}

	data := &core.AssetJournalData{}

	if err := self.assets.PutJournal(self.asset.Header(), HugActionBirth, self.receipient, data); err != nil {
		return errors.Wrap(err, "failed to update journal")
	}

	return nil
}
