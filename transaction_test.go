package chug_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crypto-hug/crypto-hug"
	"github.com/crypto-hug/crypto-hug/utils"
)

func TestLock(t *testing.T) {
	privKey := utils.GeneratePrivKey()

	tx := chug.NewTransaction(chug.SpawnGenesisHugTransactionType)
	tx.Version = "1.0.0"
	tx.Timestamp = 1562764857
	tx.Data = []byte("hug the planed")

	tx.HashTx()

	privK := utils.PrivKeyToBytes(privKey)
	pubK := utils.PubKeyToBytes(&privKey.PublicKey)

	tx.IssuerPubKey = pubK
	tx.ValidatorPubKey = pubK

	err := tx.LockIssuer(privK)
	assert.NoError(t, err)
	err = tx.LockValidator(privK)
	assert.NoError(t, err)

	err = tx.CheckLockIssuer()
	assert.NoError(t, err)

	err = tx.CheckLockValidator()
	assert.NoError(t, err)

	assert.NotEmpty(t, tx.IssuerLock)
	assert.NotEmpty(t, tx.ValidatorLock)

	assert.Equal(t, tx.ValidatorLock, tx.IssuerLock)

	fmt.Printf("priv key:   %s\n", utils.Base58ToStr(privK))
	fmt.Printf("pub key:    %s\n", utils.Base58ToStr(pubK))
	fmt.Printf("hash:       %s\n", utils.Base58ToStr(tx.Hash))
	fmt.Printf("lock-1:     %s\n", utils.Base58ToStr(tx.IssuerLock))
	fmt.Printf("lock-2:     %s\n", utils.Base58ToStr(tx.ValidatorLock))
}
