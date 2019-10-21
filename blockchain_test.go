package chug_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/v-braun/go-must"

	"github.com/crypto-hug/crypto-hug/utils"

	"github.com/stretchr/testify/assert"

	"github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/fs"
)

type txWithSecret struct {
	*chug.Transaction
	privK []byte
	pubK  []byte

	txAddrStr string
}

func (t *txWithSecret) String() string {
	return `
	version: ` + string(t.Transaction.Version) + `
	timestamp: ` + strconv.FormatInt(t.Transaction.Timestamp, 10) + `
	hash: ` + t.Transaction.Hash.String() + `
	pubKey: ` + t.Transaction.IssuerPubKey.String() + `
	address: ` + t.txAddrStr + `
	lock: ` + t.Transaction.IssuerLock.String() + `
	data: ` + t.Transaction.Data.String() + `

	privKey: ` + utils.Base58ToStr(t.privK) + `
`
}

func newTestGenesisTxWithSecret() *txWithSecret {
	result := new(txWithSecret)
	result.Transaction = chug.NewTransaction(chug.SpawnGenesisHugTransactionType)
	result.privK, result.pubK, _ = utils.CreateKeyPair()

	result.Transaction.IssuerPubKey = utils.NewBase58JsonValFromData(result.pubK)
	result.Transaction.ValidatorPubKey = utils.NewBase58JsonValFromData(result.pubK)
	result.Transaction.Data = utils.NewBase58JsonValFromData([]byte("hug the tests"))

	result.HashTx()
	must.NoError(result.LockIssuer(result.privK, result.pubK), "")
	must.NoError(result.LockValidator(result.privK, result.pubK), "")

	result.txAddrStr, _ = result.Address()

	return result
}

func newTestConfig(gen *txWithSecret) *chug.Config {
	result := chug.NewDefaultConfig()
	result.GenesisTx.Address = gen.txAddrStr
	result.GenesisTx.Data = gen.Transaction.Data.String()
	result.GenesisTx.Hash = gen.Transaction.Hash.String()
	result.GenesisTx.Lock = gen.Transaction.IssuerLock.String()
	result.GenesisTx.PubKey = gen.Transaction.IssuerPubKey.String()
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

func newTestGiveHugTx(issuer *chug.Wallet, validator *chug.Wallet, bc *chug.Blockchain) *chug.Transaction {
	issuerETag, err := bc.States().HugGetEtag(issuer.Addr)
	must.NoError(err, "could not get etag for issuer")

	validatorETag, err := bc.States().HugGetEtag(validator.Addr)
	must.NoError(err, "could not get etag for validator")

	tx := chug.NewTransaction(chug.GiveHugTransactionType)
	tx.IssuerEtag = issuerETag
	tx.ValidatorEtag = validatorETag
	tx.HashTx()

	err = tx.LockIssuer(issuer.PrivK, issuer.PubK)
	must.NoError(err, "could not lock issuer")

	err = tx.LockValidator(validator.PrivK, validator.PubK)
	must.NoError(err, "could not lock validator")

	return tx
}

func newTestBlockchain() (*chug.Blockchain, *txWithSecret) {
	genTx := newTestGenesisTxWithSecret()
	cfg := newTestConfig(genTx)
	fs := newFileFSTestDir()
	bc := chug.NewBlockchain(fs, cfg)

	err := bc.ProcessTransaction(genTx.Transaction)
	if err != nil {
		panic(fmt.Sprintf("could not process gen tx: %+v", err))
	}

	return bc, genTx
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

	k1Pr, k1Pu, _ := utils.CreateKeyPair()
	k2Pr, k2Pu, _ := utils.CreateKeyPair()

	tx := chug.NewTransaction(chug.GiveHugTransactionType)
	tx.IssuerPubKey = utils.NewBase58JsonValFromData(k1Pu)
	tx.ValidatorPubKey = utils.NewBase58JsonValFromData(k2Pu)
	tx.HashTx()
	tx.LockIssuer(k1Pr, k1Pr)
	tx.LockValidator(k2Pr, k2Pu)

	bc := chug.NewBlockchain(fs, cfg)
	err = bc.ProcessTransaction(tx)

	assert.Error(t, err, err.Error())
}

func TestBlockSizeLimit(t *testing.T) {
	fs := newFileFSTestDir()

	genTx := newTestGenesisTxWithSecret()
	cfg := newTestConfig(genTx)
	cfg.Blocks.Size = 1

	bc := chug.NewBlockchain(fs, cfg)
	bc.ProcessTransaction(genTx.Transaction)
	assert.Equal(t, 1, bc.Store().BlockCount())

	usr1 := chug.NewWalletFromKeys(genTx.privK, genTx.pubK)
	usr2 := chug.NewWallet()
	tx := newTestGiveHugTx(usr1, usr2, bc)

	bc.ProcessTransaction(tx)
	assert.Equal(t, 2, bc.Store().BlockCount())

}

func TestSelfHugNotPossible(t *testing.T) {
	bc, genTx := newTestBlockchain()

	usr1 := chug.NewWalletFromKeys(genTx.privK, genTx.pubK)

	tx := newTestGiveHugTx(usr1, usr1, bc)

	err := bc.ProcessTransaction(tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "self hugging is not possible for address")
}

func TestCreateMultipleHugs(t *testing.T) {
	// fs := fs.NewFs4Tests()

	bc, genTx := newTestBlockchain()

	usr1 := chug.NewWalletFromKeys(genTx.privK, genTx.pubK)
	usr2 := chug.NewWallet()

	tx := newTestGiveHugTx(usr1, usr2, bc)

	err := bc.ProcessTransaction(tx)
	assert.NoError(t, err)
}

func _TestGenerateGenTx(t *testing.T) {
	genTx := newTestGenesisTxWithSecret()
	cfg := newTestConfig(genTx)
	fs := newFileFSTestDir()
	bc := chug.NewBlockchain(fs, cfg)
	err := bc.ProcessTransaction(genTx.Transaction)
	assert.NoError(t, err)

	fmt.Println(`

GENERATED NEW GENESIS TX
============================
` + genTx.String() + `

____________________________
`)
}

// func TestPrivKey(t *testing.T) {
// 	k, _ := rsa.GenerateKey(rand.Reader, 2048)
// 	rsaKeyBytes := x509.MarshalPKCS1PrivateKey(k)
// 	fmt.Printf("rsa len: %d\n", len(rsaKeyBytes))

// 	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	priv := privKey.D.Bytes()
// 	//pub := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
// 	fmt.Printf("ec len: %d\n", len(priv))

// }

// // func TestCreateGenesisTx(t *testing.T) {
// // 	d, _ := os.Getwd()
// // 	fs := fs.NewFileFs(d + "/testdata/")
// // 	bc := chug.NewBlockchain(fs)
// // 	bc.CreateGenesisTxIfNotExists()
// // }
