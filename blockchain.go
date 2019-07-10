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
	if !configExists(fs) {
		err := configCreateDefault(fs)
		must.NoError(err, "config creation failed")
	}

	conf, err := configLoad(fs)
	must.NoError(err, "config loading failed")

	return conf
}

func (bc *Blockchain) createGenesisTx() {
	ghPubKey, err := utils.Base58FromString(bc.config.GenesisTx.PubKey)
	must.NoError(err, "invalid config val 'GenesisTx.PubKey' (%s). Could not be Base58 decoded", bc.config.GenesisTx.PubKey)

	addr, err := HugAddrFromPubKey(ghPubKey)
	must.NoError(err, "invalid config val 'GenesisTx.PubKey' (%s). Could not cerate address", bc.config.GenesisTx.PubKey)

}
