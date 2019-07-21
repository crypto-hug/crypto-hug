package chug

import (
	"bytes"
	"fmt"
	"time"

	"github.com/crypto-hug/crypto-hug/utils"
	"github.com/pkg/errors"
	"github.com/v-braun/go-must"
)

type TransactionType string

var GiveHugTransactionType = TransactionType("HUG")
var SpawnGenesisHugTransactionType = TransactionType("GHUG")

type Transactions []*Transaction

type Transaction struct {
	Version   Version
	Type      TransactionType
	Timestamp int64

	Hash *utils.Base58JsonVal

	IssuerPubKey *utils.Base58JsonVal
	IssuerLock   *utils.Base58JsonVal
	IssuerEtag   string

	ValidatorPubKey *utils.Base58JsonVal
	ValidatorLock   *utils.Base58JsonVal
	ValidatorEtag   string

	Data *utils.Base58JsonVal
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
	genesisTx.Data, err = utils.NewBase58JsonValFromString(config.GenesisTx.Data)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Data", config.GenesisTx.Data)
	genesisTx.Hash, err = utils.NewBase58JsonValFromString(config.GenesisTx.Hash)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Hash", config.GenesisTx.Hash)

	genesisTx.Timestamp = config.GenesisTx.Timestamp
	genesisTx.Version = Version(config.GenesisTx.Version)

	genesisTx.IssuerPubKey, err = utils.NewBase58JsonValFromString(config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", config.GenesisTx.PubKey)
	genesisTx.IssuerLock, err = utils.NewBase58JsonValFromString(config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", config.GenesisTx.Lock)

	genesisTx.ValidatorPubKey, err = utils.NewBase58JsonValFromString(config.GenesisTx.PubKey)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.PubKey", config.GenesisTx.PubKey)
	genesisTx.ValidatorLock, err = utils.NewBase58JsonValFromString(config.GenesisTx.Lock)
	must.NoError(err, "could not parse %s (%s)", "config.GenesisTx.Lock", config.GenesisTx.Lock)

	err = genesisTx.Check()
	must.NoError(err, "genesis tx failed checks: %s", err)

	return genesisTx
}

func (tx *Transaction) Check() error {
	if err := tx.CheckHash(); err != nil {
		return errors.Wrap(err, "failed hash check")
	}
	if res := tx.CheckLockIssuer(); res == false {
		return errors.New("failed issuer lock check")
	}
	if res := tx.CheckLockValidator(); res == false {
		return errors.New("failed validator lock check")
	}

	return nil
}

func (tx *Transaction) CheckHash() error {
	hash := tx.getHash()
	if !bytes.Equal(hash, tx.Hash.Bytes()) {
		return errors.New(fmt.Sprintf("stored (%s) hash and actual hash (%s) are not equal", tx.Hash, utils.Base58ToStr(hash)))
	}

	return nil
}

func (tx *Transaction) CheckLockIssuer() bool {
	res := tx.checkLock(tx.IssuerPubKey.Bytes(), tx.IssuerLock.Bytes())
	return res
}

func (tx *Transaction) CheckLockValidator() bool {
	res := tx.checkLock(tx.ValidatorPubKey.Bytes(), tx.ValidatorLock.Bytes())
	return res
}

func (tx *Transaction) LockIssuer(privKey []byte, pubKey []byte) error {
	lock, err := tx.lock(privKey, pubKey)
	if err == nil {
		tx.IssuerLock = utils.NewBase58JsonValFromData(lock)
	}

	return err
}

func (tx *Transaction) LockValidator(privKey []byte, pubKey []byte) error {
	lock, err := tx.lock(privKey, pubKey)
	if err == nil {
		tx.ValidatorLock = utils.NewBase58JsonValFromData(lock)
	}

	return err
}

func (tx *Transaction) HashTx() {
	tx.Hash = utils.NewBase58JsonValFromData(tx.getHash())
}

func (tx *Transaction) Address() (string, error) {
	result, err := NewAddress(tx.Hash.Bytes())
	return result, err
}

func (tx *Transaction) IssuerAddress() (string, error) {
	result, err := NewAddress(tx.IssuerPubKey.Bytes())
	return result, err
}

func (tx *Transaction) ValidatorAddress() (string, error) {
	result, err := NewAddress(tx.ValidatorPubKey.Bytes())
	return result, err
}

func (tx *Transaction) lock(privKey []byte, pubKey []byte) ([]byte, error) {
	if len(tx.Hash.Bytes()) == 0 {
		panic(errors.New("transaction is not hashed"))
	}

	lock, err := utils.SignCreate(privKey, pubKey, tx.Hash.Bytes())
	return lock, err
}

func (tx *Transaction) checkLock(pubKey []byte, lock []byte) bool {
	if len(tx.Hash.Bytes()) == 0 {
		panic(errors.New("transaction is not hashed"))
	}

	res := utils.SignCheck(pubKey, tx.Hash.Bytes(), lock)

	return res
}

func (tx *Transaction) IsGenesisTx(conf *Config) bool {
	if tx.Hash.String() == conf.GenesisTx.Hash &&
		tx.IssuerLock.String() == conf.GenesisTx.Lock &&
		tx.IssuerPubKey.String() == conf.GenesisTx.PubKey &&
		tx.ValidatorPubKey.String() == conf.GenesisTx.PubKey &&
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
		tx.Data.Bytes(),
	}, []byte{})

	result := utils.Hash(data)
	return result
}

func (self Transactions) getHash() []byte {
	var all [][]byte
	for _, tx := range self {
		var hash = tx.Hash
		all = append(all, hash.Bytes())
	}

	result := bytes.Join(all, []byte{})
	return result
}
