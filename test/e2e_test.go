package test

import (
	"fmt"
	"log"
	"syscall"
	"testing"
	"time"

	"github.com/brandonto/rest-api-microservice-demo/core"
	"github.com/brandonto/rest-api-microservice-demo/db"

	"github.com/stretchr/testify/suite"

	bolt "go.etcd.io/bbolt"
)

type EndToEndTestSuite struct {
	suite.Suite
	dbFilePath   string
	dbBucketName string
	dbBucketKey  []byte
	port         uint64
}

// One time set up function called when the test suite is started
//
func (suite *EndToEndTestSuite) SetupSuite() {
	suite.dbFilePath = "/home/brandonto/rest-api-microservice-demo.db"
	suite.dbBucketName = "E2ETestBucket"
	suite.dbBucketKey = []byte(suite.dbBucketName)
	suite.port = 54321
}

// One time tear down function called when the test suite is completed
//
func (suite *EndToEndTestSuite) TearDownSuite() {
	// This function stub is here to make explicit that it was considered but
	// deemed unecessary to include any logic.
	//
}

// Called before every test to set up test fixture
//
func (suite *EndToEndTestSuite) SetupTest() {
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

	// Configure and run server in a seperate goroutine
	//
	dbCfg := db.Config{
		FilePath:   suite.dbFilePath,
		BucketName: suite.dbBucketName,
	}

	coreCfg := core.Config{
		DbCfg:        dbCfg,
		Port:         suite.port,
		EnableLogger: true,
		Standalone:   false,
	}

	go core.Run(coreCfg)

	// Gives a bit of breathing room to allow the server to start up
	//
	time.Sleep(1 * time.Second)
}

// Called after every test to tear down test fixture
//
func (suite *EndToEndTestSuite) TearDownTest() {
	// SIGUSR1 is sent to a channel owned by the goroutine running the core. This
	// signal performs a graceful shutdown of the server.
	//
	err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	if err != nil {
		log.Fatal(err)
	}

	// Gives a bit of breathing room to allow the server to shut down
	//
	sleepTimeInSeconds := core.ServerShutdownTimeoutInSeconds + 1
	time.Sleep(time.Duration(sleepTimeInSeconds) * time.Second)
}

func (suite *EndToEndTestSuite) TestBasicFunctionality() {
	fmt.Println("Running TestSomething()")
}

func (suite *EndToEndTestSuite) TestSomethingElse() {
	fmt.Println("Running TestSomethingElse()")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndTestSuite))
}
