package chug

import (
	"bytes"
	"fmt"
	"time"

	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
	"github.com/v-braun/must"
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
	IssuerEtag   string

	ValidatorPubKey []byte
	ValidatorLock   []byte
	ValidatorEtag   string

	Data []byte
}

func NewTransaction(t TransactionType) *Transaction {
	tx := new(Transaction)
	tx.Version = TxVersion
	tx.Timestamp = time.Now().Unix()
	tx.Type = t

	return tx
}

func NewGenesisTransaction(config *Config) *Transaction {
	var err error
	genesisTx := NewTransaction(SpawnGenesisHugTransactionType)
	genesisTx.Data, err = utils.Base58FromString(config.GenesisTx.Data)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Data", config.GenesisTx.Data)
	genesisTx.Hash, err = utils.Base58FromString(config.GenesisTx.Hash)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Hash", config.GenesisTx.Hash)

	genesisTx.Timestamp = config.GenesisTx.Timestamp
	genesisTx.Version = Version(config.GenesisTx.Version)

	genesisTx.IssuerPubKey, err = utils.Base58FromString(config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", config.GenesisTx.PubKey)
	genesisTx.IssuerLock, err = utils.Base58FromString(config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", config.GenesisTx.Lock)

	genesisTx.ValidatorPubKey, err = utils.Base58FromString(config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", config.GenesisTx.PubKey)
	genesisTx.ValidatorLock, err = utils.Base58FromString(config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", config.GenesisTx.Lock)

	err = genesisTx.Check()
	must.NoError(err, "genesis tx failed checks")

	return genesisTx
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

func (tx *Transaction) Address() (string, error) {
	result, err := NewAddress(tx.Hash)
	return result, err
}

func (tx *Transaction) IssuerAddress() (string, error) {
	result, err := NewAddress(tx.IssuerPubKey)
	return result, err
}

func (tx *Transaction) ValidatorAddress() (string, error) {
	result, err := NewAddress(tx.ValidatorPubKey)
	return result, err
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

func (tx *Transaction) IsGenesisTx(conf *Config) bool {
	if bytes.Compare(tx.Hash, utils.Base58FromStringMust(conf.GenesisTx.Hash)) == 0 &&
		bytes.Compare(tx.IssuerLock, utils.Base58FromStringMust(conf.GenesisTx.Lock)) == 0 &&
		bytes.Compare(tx.ValidatorLock, utils.Base58FromStringMust(conf.GenesisTx.Lock)) == 0 &&
		bytes.Compare(tx.IssuerPubKey, utils.Base58FromStringMust(conf.GenesisTx.PubKey)) == 0 &&
		bytes.Compare(tx.ValidatorPubKey, utils.Base58FromStringMust(conf.GenesisTx.PubKey)) == 0 &&
		tx.Timestamp == conf.GenesisTx.Timestamp &&
		tx.Type == SpawnGenesisHugTransactionType &&
		tx.Version == Version(conf.GenesisTx.Version) {
		return true
	}

	return false
}

func (tx *Transaction) getHash() []byte {
	data := bytes.Join([][]byte{
		[]byte(tx.Version),
		[]byte(tx.Type),
		[]byte(utils.Int64GetBytes(tx.Timestamp)),
		[]byte(tx.IssuerEtag),
		[]byte(tx.ValidatorEtag),
		tx.Data,
	}, []byte{})

	result := utils.Hash(data)
	return result
}

func (self Transactions) getHash() []byte {
	var all [][]byte
	for _, tx := range self {
		var hash = tx.Hash
		all = append(all, hash)
	}

	result := bytes.Join(all, []byte{})
	return result
}
