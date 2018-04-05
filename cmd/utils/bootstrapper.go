package utils

import (
	"github.com/crypto-hug/crypto-hug/core"
	"github.com/crypto-hug/crypto-hug/core/chug"
	"github.com/crypto-hug/crypto-hug/log"
	"github.com/crypto-hug/crypto-hug/persistance/blocks"
	"github.com/crypto-hug/crypto-hug/persistance/wallets"
)

func SetupBlockchain() (*core.Blockchain, error) {
	walletsPath := "./blockhain_data/wallets"
	blocksPath := "./blockhain_data/blocks"
	statsPath := "./blockhain_data/_index.db"

	log.Global().Debug("create blockchain", log.More{"blocks": blocksPath, "blockstats": statsPath})

	blockSink := blocks.NewFsBlockSink(blocksPath)
	walletSink := wallets.NewFsWalletSink(walletsPath)

	cfg := core.NewBlockchainConfig()
	cfg.CreateGenesisTransactions = createGenesisTransactions
	cfg.TransactionProcessors = core.TransactionProcessors{&core.CommonTxProcessor{},
		chug.NewSpawnHugProcessor(walletSink)}

	stat, err := blocks.NewBoltBlockStats(statsPath)
	if err != nil {
		return nil, err
	}

	store := core.NewBlockStore(blockSink, stat)

	result := core.NewBlockchain(cfg, store)

	return result, nil
}

func createGenesisTransactions() (core.Transactions, error) {
	myAddr := "VB6QzPAL7P83N48MhoFdLXuroxPmUiphp"
	log.Global().Info("create genesis block", log.More{"reward": myAddr})
	genesisOwnerAddress, err := core.NewAddressFromString(myAddr)
	if err != nil {
		return nil, err
	}
	spawnHugs := core.Transactions{}

	for i := 0; i < 3; i++ {
		var spawnTx *core.Transaction = nil
		spawnTx, err = chug.NewSpawnHugTransaction(genesisOwnerAddress)
		if err != nil {
			return nil, err
		}
		spawnHugs = append(spawnHugs, spawnTx)
	}

	return spawnHugs, err
}
