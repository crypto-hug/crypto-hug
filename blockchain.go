package chug

import (
	"github.com/crypto-hug/crypto-hug/fs"
	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/v-braun/must"
)

type Blockchain struct {
	config *Config
	fs     *fs.FileSystem
}

func NewBlockchain(fs *fs.FileSystem) *Blockchain {
	bc := new(Blockchain)
	bc.fs = fs
	bc.config = createConf(fs)

	return bc
}

func createConf(fs *fs.FileSystem) *Config {
	err := configCreateDefaultIfNotExists(fs)
	must.NoError(err, "config creation failed")

	conf, err := configLoad(fs)
	must.NoError(err, "config loading failed")

	return conf
}

func (bc *Blockchain) CreateGenesisTxIfNotExists() {
	ghPubKey, err := utils.Base58FromString(bc.config.GenesisTx.PubKey)
	must.NoError(err, "invalid config val 'GenesisTx.PubKey' (%s). Could not be Base58 decoded", bc.config.GenesisTx.PubKey)

	addr, err := HugAddrFromPubKey(ghPubKey)
	must.NoError(err, "invalid config val 'GenesisTx.PubKey' (%s). Could not cerate address", bc.config.GenesisTx.PubKey)

	hugFilePath := bc.config.Paths.HugsDir + addr + ".hug"
	if bc.fs.FileExists(hugFilePath) {
		return
	}

	genesisTx := NewTransaction(SpawnGenesisHugTransactionType)
	genesisTx.Data, err = utils.Base58FromString(bc.config.GenesisTx.Data)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Data", bc.config.GenesisTx.Data)
	genesisTx.Hash, err = utils.Base58FromString(bc.config.GenesisTx.Hash)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Hash", bc.config.GenesisTx.Hash)

	genesisTx.Timestamp = bc.config.GenesisTx.Timestamp
	genesisTx.Version = Version(bc.config.GenesisTx.Version)

	genesisTx.IssuerPubKey, err = utils.Base58FromString(bc.config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", bc.config.GenesisTx.PubKey)
	genesisTx.IssuerLock, err = utils.Base58FromString(bc.config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", bc.config.GenesisTx.Lock)

	genesisTx.ValidatorPubKey, err = utils.Base58FromString(bc.config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", bc.config.GenesisTx.PubKey)
	genesisTx.ValidatorLock, err = utils.Base58FromString(bc.config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", bc.config.GenesisTx.Lock)

	err = genesisTx.Check()
	must.NoError(err, "genesis tx failed checks")

	jsonData, err := utils.JsonSerializeRaw(genesisTx)
	must.NoError(err, "genesis tx failed serialization")

	err = bc.fs.WriteFile(hugFilePath, jsonData)
	must.NoError(err, "genesis tx failed write to file")
}
