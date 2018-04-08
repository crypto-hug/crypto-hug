package utils

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/core/chug"
	"github.com/crypto-hug/crypto-hug/formatters"
	"github.com/crypto-hug/crypto-hug/log"
	"github.com/crypto-hug/crypto-hug/persistance/assets"
	"github.com/crypto-hug/crypto-hug/persistance/blocks"
	"github.com/crypto-hug/crypto-hug/persistance/wallets"
)

func SetupBlockchain() (*core.Blockchain, error) {
	assetsPath := "./blockhain_data/assets"
	walletsPath := "./blockhain_data/wallets"
	blocksPath := "./blockhain_data/blocks"
	statsPath := "./blockhain_data/_index.db"

	log.Global().Debug("create blockchain", log.More{"blocks": blocksPath, "blockstats": statsPath})

	blockSink := blocks.NewFsBlockSink(blocksPath)
	walletSink := wallets.NewFsWalletSink(walletsPath)
	assetSink := assets.NewFsAssetSink(assetsPath)

	stat, err := blocks.NewBoltBlockStats(statsPath)
	if err != nil {
		return nil, err
	}
	store := core.NewBlockStore(blockSink, stat)

	cfg := core.NewBlockchainConfig()
	cfg.CreateGenesisTransactions = createGenesisTransactions

	cfg.CreateTxProcessors = createTxProcessors

	if err != nil {
		return nil, err
	}

	result := core.NewBlockchain(cfg, store, walletSink, assetSink)

	return result, nil
}

func createTxProcessors() core.TransactionProcessors {
	result := core.TransactionProcessors{&core.CommonTxProcessor{},
		&chug.SpawnHugProcessor{}, &chug.SpendHugProcessor{}}

	return result
}

func createGenesisTransactions() (core.Transactions, error) {
	myAddr := "VB6QzPAL7P83N48MhoFdLXuroxPmUiphp"
	myPubKey, err := formatters.Base58FromString("aY3JXGjbhvc8gpyepFpGgEhoKmjLYL8piWeKP6cYNGP9U5zs9HqHJASBb7WbD5FevKTJWvhcctZd5w3Em62bquoVB6QzPAL7P83N48MhoFdLXuroxPmUiphp")
	if err != nil {
		panic(err)
	}

	log.Global().Info("create genesis block", log.More{"reward": myAddr})
	genesisOwnerAddress, err := core.NewAddressFromString(myAddr)
	if err != nil {
		return nil, err
	}
	spawnHugs := core.Transactions{}

	for i := 0; i < 3; i++ {
		var spawnTx *core.Transaction = nil
		spawnTx, err = chug.NewSpawnHugTransaction(genesisOwnerAddress, myPubKey)
		if err != nil {
			return nil, err
		}
		spawnHugs = append(spawnHugs, spawnTx)
	}

	return spawnHugs, err
}
