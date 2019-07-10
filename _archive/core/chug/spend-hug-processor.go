package chug

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/log"
)

type SpendHugProcessor struct {
	wallets core.WalletSink
	assets  core.AssetSink

	sender     *core.Address
	receipient *core.Address

	asset *core.AssetHeader
}

func (self *SpendHugProcessor) Setup(wallets core.WalletSink, assets core.AssetSink) {
	errors.AssertNotNil("wallets", wallets)
	errors.AssertNotNil("assets", assets)

	self.wallets = wallets
	self.assets = assets
}
func (self *SpendHugProcessor) ShouldProcess(tx *core.Transaction) bool {
	return tx.Type == SpendHugTxType
}

func (self *SpendHugProcessor) Name() string {
	return "SpendHug"
}

func (self *SpendHugProcessor) Prepare(tx *core.Transaction) error {
	logger := log.NewLog("SpendHugProcessor")

	errors.AssertNotNil("tx", tx)
	errors.AssertNotNil("self.wallets", self.wallets)
	errors.AssertNotNil("self.assets", self.assets)
	errors.AssertTrue("tx.Type", tx.Type == SpendHugTxType, "is not SpendHug")

	sender := core.NewAddressFromStringStrict(tx.Sender)

	data, err := UnwrapSpendHugTxData(tx)
	if err != nil {
		return errors.Wrap(err, "invalid SpendHug tx data")
	}

	receipient := core.NewAddressFromStringStrict(data.RecipientAddress)
	assetAddress := core.NewAddressFromStringStrict(data.HugAddress)

	asset, err := self.assets.GetHeader(assetAddress)
	if err != nil {
		if err == core.AssetNotExist {
			logger.Error("could not find asset", log.More{"address": assetAddress.Address})
		}

		return err
	}

	senderOwnsAsset, err := self.wallets.HasAsset(sender, asset)
	if err != nil {
		return errors.Wrap(err, "could not get sender owned assets")
	}
	if senderOwnsAsset == false {
		return core.InsufficientFunds
	}

	receipientOwnedAssets, err := self.wallets.ListAssetsByType(receipient, AssetTypeHug)
	if err != nil {
		return errors.Wrap(err, "could not get owned assets")
	}

	if receipientOwnedAssets != nil && len(receipientOwnedAssets) >= max_hugs_in_wallet {
		return errors.NewErrorFromString("max hug balance of %v reached: current %v", max_hugs_in_wallet, len(receipientOwnedAssets))
	}

	self.asset = asset
	self.receipient = receipient
	self.sender = sender

	return nil
}

func (self *SpendHugProcessor) Process(tx *core.Transaction) error {
	if err := self.wallets.RemoveAsset(self.sender, self.asset); err != nil {
		return errors.Wrap(err, "failed to set new balance")
	}
	if err := self.wallets.PutAssetPropF(self.receipient, self.asset, "balance", 1); err != nil {
		return errors.Wrap(err, "failed to set new balance")
	}

	data := &core.AssetJournalData{"sender": self.sender.Address,
		"receipient": self.receipient.Address}
	if err := self.assets.PutJournal(self.asset, HugActionSpend, self.sender, data); err != nil {
		return errors.Wrap(err, "failed to update journal")
	}

	return nil
}
