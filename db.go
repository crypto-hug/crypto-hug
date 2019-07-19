package chug

import (
	"fmt"
	"sort"

	"github.com/v-braun/must"
	"go.etcd.io/bbolt"
)

type bucketName []byte

var bucketHugLink = bucketName("hug_link")

type dbTx struct {
	*bbolt.Tx
}

type DB struct {
	*bbolt.DB
}

func NewDB(path string) *DB {
	db, err := bbolt.Open(path, 0666, nil)
	must.NoError(err, "could not opend db (%s)", path)

	result := new(DB)
	result.DB = db

	return result
}

func (db *DB) NewWritableTx() (*dbTx, error) {
	result := new(dbTx)
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	result.Tx = tx
	return result, nil
}

func (tx *dbTx) hugLinkCreate(issuer string, validator string, transaction string) bool {
	addresses := []string{issuer, validator}
	sort.Strings(addresses)

	link := fmt.Sprintf("%s-%s", addresses[0], addresses[1])

	bucket := tx.Bucket(bucketHugLink)
	val := bucket.Get([]byte(link))
	if val == nil {
		return false
	}

	bucket.Put([]byte(link), []byte(transaction))
	return true
}
