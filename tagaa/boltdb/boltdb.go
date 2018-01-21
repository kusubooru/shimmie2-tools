package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kusubooru/tagaa/tagaa"
)

const (
	imageBucket = "images"
	groupBucket = "groups"
	blobBucket  = "blobs"
)

type store struct {
	*bolt.DB
}

func (db *store) Close() error {
	return db.DB.Close()
}

// openBolt creates and opens a bolt database at the given path. If the file does
// not exist then it will be created automatically. After opening it creates
// all the needed buckets.
func openBolt(file string) (*bolt.DB, error) {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening bolt file: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(imageBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(groupBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(blobBucket))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("creating buckets: %v", err)
	}
	return db, nil
}

// NewStore opens the bolt database file and returns an implementation for
// tagaa.Store. The bolt database file will be created if it does not exist.
func NewStore(boltFile string) (tagaa.Store, error) {
	db, err := openBolt(boltFile)
	if err != nil {
		return nil, err
	}
	return &store{db}, nil
}