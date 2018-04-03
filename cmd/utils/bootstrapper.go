package utils

import (
	"../../core"
	"../../core/txvalidators"
	"../../log"
	"../../persistance/blocks"
	//ishell "gopkg.in/abiosoft/ishell.v2"
)

func SetupBlockchain() (*core.Blockchain, error) {
	blocksPath := "./blockhain_data/blocks"
	statsPath := "./blockhain_data/_index.db"

	log.Global().Debug("create blockchain", log.More{"blocks": blocksPath, "blockstats": statsPath})

	sink := blocks.NewFsBlockSink(blocksPath)
	stat, err := blocks.NewBoltBlockStats(statsPath)
	if err != nil {
		return nil, err
	}

	store := core.NewBlockStore(sink, stat)

	txvReg := txvalidators.SharedRegistry()

	result := core.NewBlockchain(store, txvReg)

	return result, nil
}

// func SetupDb() *hugdb.BoltDb {
// 	var filePath = "./blockhain_data/c_hug.db"
// 	log.Global().Debug("use blockchain file", log.More{"file": filePath})

// 	var db, err = hugdb.NewBoltDB(filePath)
// 	if err != nil {
// 		FatalExit(err)
// 		return nil
// 	}

// 	db.Bootstrap()

// 	return db
// }
