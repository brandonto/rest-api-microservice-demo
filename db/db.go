package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/brandonto/rest-api-microservice-demo/model"

	bolt "go.etcd.io/bbolt"
)

// Structure to abtract away the underlying database implementation. This is the
// main handle into the database.
//
type Db struct {
	boltDb    *bolt.DB
	bucketKey []byte
	Config
}

// Structure to encapsulate configuration needed to initialize the Db
//
type Config struct {
	FilePath   string
	BucketName string
}

// Constructor for Db object
//
func NewDb(config Config) *Db {
	return &Db{Config: config}
}

// Creates/opens and initializes a bbolt DB
//
func (db *Db) Initialize() error {
	boltDb, err := bolt.Open(db.FilePath, 0600, nil)
	if err != nil {
		return err
	}

	db.boltDb = boltDb
	db.bucketKey = []byte(db.BucketName)

	err = db.boltDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(db.bucketKey)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Closes the Db. Not strictly necessary in this application, but good practice
// regardless
//
func (db *Db) Close() {
	db.boltDb.Close()
}

// Returns a list of up to "limit" number of DetailedMessage starting with the
// first entry from index "id". A non-0 "afterId" returned indicates that there
// are more messages left to retrieve from the database.
//
func (db *Db) ListMessages(limit uint64, id uint64) ([]*model.DetailedMessage, uint64, error) {
	var detailedMessages []*model.DetailedMessage
	afterId := uint64(0)

	err := db.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucketKey)
		if bucket == nil {
			// This shouldn't be possible as the bucket should have been created
			// on initialization. Fundamental flaw in program operation... so
			// lets just die
			//
			log.Fatal(errors.New("Irrecoverable state"))
		}

		cursor := bucket.Cursor()
		numMessagesRetrieved := uint64(0)
		for k, v := cursor.Seek(uint64ToBytes(id)); k != nil; k, v = cursor.Next() {
			detailedMessage := &model.DetailedMessage{}

			if v == nil {
				// This shouldn't be possible. A nil value for a non-nil key is
				// only returned from Seek() if the key is associated with a
				// bucket value. This application should never store a bucket
				// here as only message blobs should have been stored. This is a
				// fundamental flaw in program operation... so lets just die.
				//
				log.Fatal(errors.New("Irrecoverable state"))
			}

			// Converts message data blob into application data structure
			//
			err := json.Unmarshal(v, detailedMessage)
			if err != nil {
				return err
			}

			// Add message to list of returned messages
			//
			detailedMessages = append(detailedMessages, detailedMessage)

			// We've reached our limit for messages. Do one quick peek at the
			// next message to see if any more exists. If one does exist, we'll
			// set the next "afterId" to use.
			//
			numMessagesRetrieved += 1
			if numMessagesRetrieved == limit {
				testKey, _ := cursor.Next()
				if testKey != nil {
					afterId = detailedMessage.Message.Id
				}
				break
			}
		}

		return nil
	})

	return detailedMessages, afterId, err
}

// Inserts a new DetailedMessage in the database using "detailedMessage" at the
// index specified in the message. Returns an error if something went wrong
// during the transaction.
//
func (db *Db) CreateMessage(detailedMessage *model.DetailedMessage) error {
	return db.boltDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucketKey)
		if bucket == nil {
			// This shouldn't be possible as the bucket should have been created
			// on initialization. Fundamental flaw in program operation... so
			// lets just die
			//
			log.Fatal(errors.New("Irrecoverable state"))
		}

		// Get the next unique integer identifier from the database to use as
		// the message ID and database key
		//
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		detailedMessage.Message.Id = id

		// Converts application data structure into message data blob
		//
		buf, err := json.Marshal(detailedMessage)
		if err != nil {
			return err
		}

		// Persists message data blob to database
		//
		return bucket.Put(uint64ToBytes(id), buf)
	})
}

// Retrieves a DetailedMessage from the database at index "id". Returns an error
// if something went wrong during the transaction.
//
func (db *Db) GetMessage(id uint64) (*model.DetailedMessage, error) {
	detailedMessage := &model.DetailedMessage{}

	err := db.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucketKey)
		if bucket == nil {
			// This shouldn't be possible as the bucket should have been created
			// on initialization. Fundamental flaw in program operation... so
			// lets just die
			//
			log.Fatal(errors.New("Irrecoverable state"))
		}

		// Retrieves message data blob from database
		//
		buf := bucket.Get(uint64ToBytes(id))
		if buf == nil {
			return fmt.Errorf("Unable to retrieve message (id=%d) from database.", id)
		}

		// Converts message data blob into application data structure
		//
		err := json.Unmarshal(buf, detailedMessage)
		if err != nil {
			return err
		}

		return nil
	})

	return detailedMessage, err
}

// Replaces a DetailedMessage in the database with "detailedMessage" at the
// index specified in the message. Returns an error if something went wrong
// during the transaction.
//
func (db *Db) UpdateMessage(detailedMessage *model.DetailedMessage) error {
	return db.boltDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucketKey)
		if bucket == nil {
			// This shouldn't be possible as the bucket should have been created
			// on initialization. Fundamental flaw in program operation... so
			// lets just die
			//
			log.Fatal(errors.New("Irrecoverable state"))
		}

		// Converts application data structure into message data blob
		//
		buf, err := json.Marshal(detailedMessage)
		if err != nil {
			return err
		}

		// Persists message data blob to database
		//
		id := detailedMessage.Message.Id
		return bucket.Put(uint64ToBytes(id), buf)
	})
}

// Deletes a DetailedMessage from the database at index "id". Returns an error
// if something went wrong during the transaction.
//
func (db *Db) DeleteMessage(id uint64) error {
	return db.boltDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(db.bucketKey)
		if bucket == nil {
			// This shouldn't be possible as the bucket should have been created
			// on initialization. Fundamental flaw in program operation... so
			// lets just die
			//
			log.Fatal(errors.New("Irrecoverable state"))
		}

		// Deletes message data blob from database
		//
		err := bucket.Delete(uint64ToBytes(id))
		if err != nil {
			return fmt.Errorf("Unable to delete message (id=%d) from database.", id)
		}

		return nil
	})
}

// Delete all Messages from the database. Fast way of doing so is to just delete
// and re create the bucket.
//
// This function isn't called in the core application. It's created to be called
// by the unit testing code.
//
func (db *Db) ClearMessages() error {
	return db.boltDb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(db.bucketKey)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}

		_, err = tx.CreateBucket(db.bucketKey)
		if err != nil {
			return err
		}

		return nil
	})
}
