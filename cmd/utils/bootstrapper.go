package utils

import (
	"../../core"
	"../../core/chug"
	"../../core/txvalidators"
	"../../log"
	"../../persistance/blocks"
	//ishell "gopkg.in/abiosoft/ishell.v2"
)

func SetupBlockchain() (*core.Blockchain, error) {
	blocksPath := "./blockhain_data/blocks"
	statsPath := "./blockhain_data/_index.db"

	log.Global().Debug("create blockchain", log.More{"blocks": blocksPath, "blockstats": statsPath})

	cfg := core.NewBlockchainConfig(createGenesisTransactions)
	sink := blocks.NewFsBlockSink(blocksPath)
	stat, err := blocks.NewBoltBlockStats(statsPath)
	if err != nil {
		return nil, err
	}

	store := core.NewBlockStore(sink, stat)

	txvReg := txvalidators.SharedRegistry()

	result := core.NewBlockchain(cfg, store, txvReg)

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
