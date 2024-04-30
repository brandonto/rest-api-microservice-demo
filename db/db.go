package db

import (
    "log"

    "github.com/brandonto/rest-api-microservice-demo/model"

    bolt "go.etcd.io/bbolt"
)

// Structure to abtract away the underlying database implementation
//
type Db struct {
    boltDb *bolt.DB
}

// Creates/opens a bbolt DB at specified filePath
//
func Open(filePath string) *Db {
    db, err := bolt.Open(filePath, 0600, nil)
    if err != nil {
        log.Fatal(err)
    }

    return &Db{db}
}

// Closes a bbolt DB. Not strictly necessary in this application, but good
// practice regardless
//
func (db *Db) Close() {
    db.boltDb.Close()
}

func (db *Db) CreateMessage(message *model.DetailedMessage) {
    // TODO
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
