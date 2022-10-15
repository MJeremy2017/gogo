package data

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

type BoltDB struct {
	db         *bolt.DB
	bucketName []byte
}

// NewBoltDB creates a new default bolt DB with a bucket
func NewBoltDB(bucketName string) *BoltDB {
	db, err := bolt.Open("local.db", 0666, nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = createBucket(db, bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	return &BoltDB{db, []byte(bucketName)}
}

func createBucket(db *bolt.DB, name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *BoltDB) Set(key []byte, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)
		err := bucket.Put(key, value)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to write %s, %s. %v", string(key), string(value), err)
	}
	return nil
}

func (b *BoltDB) Get(key []byte) string {
	var value []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucketName)
		v := bucket.Get(key)
		if v == nil {
			log.Printf("key %s does not exist", string(key))
			return nil
		}
		value = v
		return nil
	})
	if err != nil {
		return ""
	}
	return string(value)
}
