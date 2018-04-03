package boltdb

import (
	// "../../serialization"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type BoltDB struct {
	db *bolt.DB
}

func NewDB(dbFile string) (*BoltDB, error) {
	dbFile, err := filepath.Abs(dbFile)
	if err != nil {
		return nil, errors.Wrap(err, "NewDB")
	}
	dir := filepath.Dir(dbFile)
	os.MkdirAll(dir, 0700)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewDB:OpenDb")
	}

	result := BoltDB{db: db}
	return &result, nil
}

func (self *BoltDB) Get(bucket string, key string) ([]byte, error) {
	var result []byte = nil

	var err = self.db.Update(func(tx *bolt.Tx) error {
		hsh, err := self.GetFromTx(tx, bucket, key)
		result = hsh
		return err
	})

	return result, errors.Wrap(err, "boltdb:Get")
}

func (self *BoltDB) GetFromTx(tx *bolt.Tx, bucket string, key string) ([]byte, error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, errors.New(fmt.Sprintf("unknown bucket %s", bucket))
	}

	var result = b.Get([]byte(key))
	return result, nil
}

func (self *BoltDB) Put(bucket string, key string, data []byte) error {
	var err = self.db.Update(func(tx *bolt.Tx) error {
		err := self.PutToTx(tx, bucket, key, data)
		return err
	})

	return errors.Wrap(err, "boltdb:Put")
}
func (self *BoltDB) PutToTx(tx *bolt.Tx, bucket string, key string, data []byte) error {
	b := tx.Bucket([]byte(bucket))
	err := b.Put([]byte(key), data)
	return errors.Wrap(err, "boltdb:PutToTx")
}

func (self *BoltDB) CreateBuckets(names ...string) error {
	var err = self.db.Update(func(tx *bolt.Tx) error {
		for _, n := range names {
			var _, cErr = tx.CreateBucketIfNotExists([]byte(n))

			if cErr != nil {
				return cErr
			}
		}

		return nil
	})

	return err
}
