package chug_test

import (
	"crypto/rsa"
	"os"
	"testing"

	"github.com/v-braun/go-must"

	"github.com/crypto-hug/crypto-hug/utils"

	"github.com/stretchr/testify/assert"

	"github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/fs"
)

type txWithSecret struct {
	*chug.Transaction
	pk *rsa.PrivateKey

	pubKeyStr  string
	privKeyStr string
	lockStr    string
	txHashStr  string
	txAddrStr  string
}

func newTestGenesisTxWithSecret() *txWithSecret {
	result := new(txWithSecret)
	result.Transaction = chug.NewTransaction(chug.SpawnGenesisHugTransactionType)
	result.pk = utils.GeneratePrivKey()

	result.Transaction.IssuerPubKey = utils.NewBase58JsonValFromData(utils.PubKeyToBytes(&result.pk.PublicKey))
	result.Transaction.ValidatorPubKey = utils.NewBase58JsonValFromData(utils.PubKeyToBytes(&result.pk.PublicKey))
	// result.Transaction.Data = utils.NewBase58JsonValFromData([]byte("hug the tests"))

	result.HashTx()
	must.NoError(result.LockIssuer(utils.PrivKeyToBytes(result.pk)), "")
	must.NoError(result.LockValidator(utils.PrivKeyToBytes(result.pk)), "")

	result.pubKeyStr = result.Transaction.IssuerPubKey.String()
	result.privKeyStr = utils.NewBase58JsonValFromData(utils.PrivKeyToBytes(result.pk)).String()
	result.lockStr = result.Transaction.IssuerLock.String()
	result.txHashStr = result.Transaction.Hash.String()
	result.txAddrStr, _ = result.Address()

	return result
}

func newTestConfig(gen *txWithSecret) *chug.Config {
	result := chug.NewDefaultConfig()
	result.GenesisTx.Address = gen.txAddrStr
	result.GenesisTx.Data = gen.Transaction.Data.String()
	result.GenesisTx.Hash = gen.txHashStr
	result.GenesisTx.Lock = gen.lockStr
	result.GenesisTx.PubKey = gen.pubKeyStr
	result.GenesisTx.Timestamp = gen.Timestamp
	result.GenesisTx.Version = string(gen.Version)
	return result
}

func newFileFSTestDir() *fs.FileSystem {
	d, _ := os.Getwd()
	fs := fs.NewFileFs(d + "/testdata/")
	fs.RemoveAll("./")
	return fs
}

func TestCreateGenesisTx(t *testing.T) {
	fs := newFileFSTestDir()
	// fs := fs.NewFs4Tests()

	cfg, err := chug.NewConfigFromFileOrDefault(fs)
	assert.NoError(t, err)

	bc := chug.NewBlockchain(fs, cfg)
	bc.CreateGenesisBlockIfNotExists()
}

func TestProcessTxWithoutGenesisTxShouldFail(t *testing.T) {
	// fs := newFileFSTestDir()
	fs := fs.NewFs4Tests()

	cfg, err := chug.NewConfigFromFileOrDefault(fs)
	assert.NoError(t, err)

	k1 := utils.GeneratePrivKey()
	k2 := utils.GeneratePrivKey()

	tx := chug.NewTransaction(chug.GiveHugTransactionType)
	tx.IssuerPubKey = utils.NewBase58JsonValFromData(utils.PubKeyToBytes(&k1.PublicKey))
	tx.ValidatorPubKey = utils.NewBase58JsonValFromData(utils.PubKeyToBytes(&k2.PublicKey))
	tx.HashTx()
	tx.LockIssuer(utils.PrivKeyToBytes(k1))
	tx.LockValidator(utils.PrivKeyToBytes(k2))

	bc := chug.NewBlockchain(fs, cfg)
	err = bc.ProcessTransaction(tx)

	assert.Error(t, err, err.Error())
}

func TestTest(t *testing.T) {
	genTx := newTestGenesisTxWithSecret()
	cfg := newTestConfig(genTx)
	fs := newFileFSTestDir()
	bc := chug.NewBlockchain(fs, cfg)
	err := bc.ProcessTransaction(genTx.Transaction)
	assert.NoError(t, err)

}

// func TestCreateGenesisTx(t *testing.T) {
// 	d, _ := os.Getwd()
// 	fs := fs.NewFileFs(d + "/testdata/")
// 	bc := chug.NewBlockchain(fs)
// 	bc.CreateGenesisTxIfNotExists()
// }
