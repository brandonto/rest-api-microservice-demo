package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/brandonto/rest-api-microservice-demo/core"
	"github.com/brandonto/rest-api-microservice-demo/db"

	"github.com/stretchr/testify/suite"

	bolt "go.etcd.io/bbolt"
)

type EndToEndTestSuite struct {
	suite.Suite
	dbFilePath  string
	dbBucketKey []byte
	port        uint64
}

func (suite *EndToEndTestSuite) SetupSuite() {
	// Configure suite
	//
	suite.dbFilePath = "/home/brandonto/rest-api-microservice-demo.db"
	dbBucketName := "E2ETestBucket"
	suite.dbBucketKey = []byte(dbBucketName)
	suite.port = 54321

	// Purges any previous data left in the bucket
	//
	boltDb, err := bolt.Open(suite.dbFilePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = boltDb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(suite.dbBucketKey)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	boltDb.Close()

	// Run server in a seperate goroutine
	//
	dbCfg := db.Config{suite.dbFilePath, dbBucketName}
	coreCfg := core.Config{dbCfg, suite.port, true}
	go core.Run(coreCfg)

	fmt.Println("[TEST] Server started.")
	time.Sleep(10 * time.Second)
}

func (suite *EndToEndTestSuite) TestSomething() {
	fmt.Println("Running TestSomething()")
}

func (suite *EndToEndTestSuite) TestSomethingElse() {
	fmt.Println("Running TestSomethingElse()")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndTestSuite))
}
