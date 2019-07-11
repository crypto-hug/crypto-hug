package chug

import (
	"bytes"
	"fmt"
	"time"

	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
)

type TransactionType string

var GiveHugTransactionType = TransactionType("HUG")
var SpawnGenesisHugTransactionType = TransactionType("GHUG")

type Transactions []*Transaction

type Transaction struct {
	Version   Version
	Type      TransactionType
	Timestamp int64

	Hash []byte

	IssuerPubKey []byte
	IssuerLock   []byte

	ValidatorPubKey []byte
	ValidatorLock   []byte

	Data []byte
}

func (tx *Transaction) Check() error {
	if err := tx.CheckHash(); err != nil {
		return errors.Wrap(err, "failed hash check")
	}
	if err := tx.CheckLockIssuer(); err != nil {
		return errors.Wrap(err, "failed issuer lock check")
	}
	if err := tx.CheckLockValidator(); err != nil {
		return errors.Wrap(err, "failed validator lock check")
	}

	return nil
}

func (tx *Transaction) CheckHash() error {
	hash := tx.getHash()
	if !bytes.Equal(hash, tx.Hash) {
		return errors.New(fmt.Sprintf("stored (%s) hash and actual hash (%s) are not equal", utils.Base58ToStr(tx.Hash), utils.Base58ToStr(hash)))
	}

	return nil
}

func (tx *Transaction) CheckLockIssuer() error {
	err := tx.checkLock(tx.IssuerPubKey, tx.IssuerLock)
	return err
}

func (tx *Transaction) CheckLockValidator() error {
	err := tx.checkLock(tx.ValidatorPubKey, tx.ValidatorLock)
	return err
}

func (tx *Transaction) LockIssuer(privKey []byte) error {
	lock, err := tx.lock(privKey)
	if err == nil {
		tx.IssuerLock = lock
	}

	return err
}

func (tx *Transaction) LockValidator(privKey []byte) error {
	lock, err := tx.lock(privKey)
	if err == nil {
		tx.ValidatorLock = lock
	}

	return err
}

func (tx *Transaction) HashTx() {
	tx.Hash = tx.getHash()
}

func NewTransaction(t TransactionType) *Transaction {
	tx := new(Transaction)
	tx.Version = TxVersion
	tx.Timestamp = time.Now().Unix()
	tx.Type = t

	return tx
}

func (tx *Transaction) lock(privKey []byte) ([]byte, error) {
	if len(tx.Hash) == 0 {
		panic(errors.New("transaction is not hashed"))
	}

	lock, err := utils.SignCreate(privKey, tx.Hash)
	return lock, err
}

func (tx *Transaction) checkLock(pubKey []byte, lock []byte) error {
	if len(tx.Hash) == 0 {
		panic(errors.New("transaction is not hashed"))
	}

	err := utils.SignCheck(pubKey, tx.Hash, lock)
	return err
}

func (tx *Transaction) getHash() []byte {
	data := bytes.Join([][]byte{
		[]byte(tx.Version),
		[]byte(tx.Type),
		[]byte(utils.Int64GetBytes(tx.Timestamp)),
		tx.Data,
	}, []byte{})

	result := utils.Hash(data)
	return result
}
