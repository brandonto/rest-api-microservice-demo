package db

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"

    "github.com/brandonto/rest-api-microservice-demo/model"

    bolt "go.etcd.io/bbolt"
)

const bucketName = "DetailedMessageBucket"

// Structure to abtract away the underlying database implementation
//
type Db struct {
    boltDb    *bolt.DB
}

// Creates/opens a bbolt DB at specified filePath and initializes it
//
// TODO fix comments
//
func Open(filePath string) *Db {
    boltDb, err := bolt.Open(filePath, 0600, nil)
    if err != nil {
        log.Fatal(err)
    }

    db := &Db{boltDb}
    db.init()

    return db
}

func (db *Db) init() {
    err := db.boltDb.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(bucketName))
        if err != nil {
            return err
        }

        // TODO more initialization may be needed

        return nil
    })

    if err != nil {
        log.Fatal(err)
    }
}

// Closes a bbolt DB. Not strictly necessary in this application, but good
// practice regardless
//
func (db *Db) Close() {
    db.boltDb.Close()
}

func (db *Db) CreateMessage(detailedMessage *model.DetailedMessage) error {
    return db.boltDb.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(bucketName))
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

func (db *Db) GetMessage(id uint64) (*model.DetailedMessage, error) {
    detailedMessage := &model.DetailedMessage{}
    err := db.boltDb.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(bucketName))
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

func (db *Db) UpdateMessage(message *model.DetailedMessage, id uint64) {
    //return db.boltDb.Update(func(tx *bolt.Tx) error {
    //    bucket := tx.Bucket([]byte(bucketName))
    //    if bucket == nil {
    //        // This shouldn't be possible as the bucket should have been created
    //        // on initialization. Fundamental flaw in program operation... so
    //        // lets just die
    //        //
    //        log.Fatal(errors.New("Irrecoverable state"))
    //    }
    //})
}

func (db *Db) DeleteMessage(id uint64) error {
    return db.boltDb.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte(bucketName))
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
