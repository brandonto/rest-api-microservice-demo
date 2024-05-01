package db

import (
    "encoding/json"
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

        id, _ := bucket.NextSequence()

        message := detailedMessage.Message
        message.Id = id

        buf, err := json.Marshal(detailedMessage)
        if err != nil {
            return err
        }

        fmt.Println(detailedMessage)
        fmt.Println(buf)

        return bucket.Put(uint64ToBytes(id), buf)
    })
}

func (db *Db) GetMessage(id int64) *model.DetailedMessage {
    // TODO
    return &model.DetailedMessage{nil, nil}
}

func (db *Db) UpdateMessage(message *model.DetailedMessage) {
    // TODO
}

func (db *Db) DeleteMessage(id int64) {
    // TODO
}
