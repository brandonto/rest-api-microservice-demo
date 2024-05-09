package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/brandonto/rest-api-microservice-demo/api"
	"github.com/brandonto/rest-api-microservice-demo/core"
	"github.com/brandonto/rest-api-microservice-demo/db"
	"github.com/brandonto/rest-api-microservice-demo/model"

	"github.com/stretchr/testify/assert"
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
	suite.dbFilePath = "/home/brandonto/rest-api-microservice-demo-test.db"
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

func (suite *EndToEndTestSuite) makeRequestURL(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", suite.port, path)
}

func (suite *EndToEndTestSuite) TestBasicFunctionality() {
	var response *http.Response
	var err error
	var message *model.Message
	var buf []byte
	var listMessagesResponse api.ListMessagesResponse
	var getMessageResponse api.GetMessageResponse
	var request *http.Request
	client := &http.Client{}

	response, err = http.Get(suite.makeRequestURL("/messages"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")

	response, err = http.Get(suite.makeRequestURL("/messages/1"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNotFound, "Unexpected HTTP status code")

	message = &model.Message{
		Payload: "foo",
	}
	buf, err = json.Marshal(message)
	assert.Nil(suite.T(), err, "Error encoding json")
	response, err = http.Post(suite.makeRequestURL("/messages"), "application/json", bytes.NewReader(buf))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusCreated, "Unexpected HTTP status code")

	response, err = http.Get(suite.makeRequestURL("/messages"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &listMessagesResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.Equal(suite.T(), len(listMessagesResponse), 1, "Unexpected number of messages in response")
	assert.Equal(suite.T(), listMessagesResponse[0].Id, uint64(1), "Unexpected message id in response")
	assert.Equal(suite.T(), listMessagesResponse[0].Payload, "foo", "Unexpected message payload response")

	response, err = http.Get(suite.makeRequestURL("/messages/1"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &getMessageResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.NotNil(suite.T(), getMessageResponse.Message, "Response did not include valid message")
	assert.Equal(suite.T(), getMessageResponse.Id, uint64(1), "Unexpected message id in response")
	assert.Equal(suite.T(), getMessageResponse.Payload, "foo", "Unexpected message payload response")

	message.Payload = "bar"
	buf, err = json.Marshal(message)
	assert.Nil(suite.T(), err, "Error encoding json")
	request, err = http.NewRequest(http.MethodPut, suite.makeRequestURL("/messages/1"), bytes.NewReader(buf))
	assert.Nil(suite.T(), err, "Error creating HTTP request")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNoContent, "Unexpected HTTP status code")

	response, err = http.Get(suite.makeRequestURL("/messages/1"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &getMessageResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.NotNil(suite.T(), getMessageResponse.Message, "Response did not include valid message")
	assert.Equal(suite.T(), getMessageResponse.Id, uint64(1), "Unexpected message id in response")
	assert.Equal(suite.T(), getMessageResponse.Payload, "bar", "Unexpected message payload response")

	request, err = http.NewRequest(http.MethodDelete, suite.makeRequestURL("/messages/1"), nil)
	assert.Nil(suite.T(), err, "Error creating HTTP request")
	response, err = client.Do(request)
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNoContent, "Unexpected HTTP status code")

	response, err = http.Get(suite.makeRequestURL("/messages/1"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNotFound, "Unexpected HTTP status code")

	message.Payload = "foo"
	buf, err = json.Marshal(message)
	assert.Nil(suite.T(), err, "Error encoding json")
	response, err = http.Post(suite.makeRequestURL("/messages"), "application/json", bytes.NewReader(buf))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusCreated, "Unexpected HTTP status code")

	message.Payload = "bar"
	buf, err = json.Marshal(message)
	assert.Nil(suite.T(), err, "Error encoding json")
	request, err = http.NewRequest(http.MethodPut, suite.makeRequestURL("/messages/1"), bytes.NewReader(buf))
	assert.Nil(suite.T(), err, "Error creating HTTP request")
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNotFound, "Unexpected HTTP status code")

	request, err = http.NewRequest(http.MethodDelete, suite.makeRequestURL("/messages/1"), nil)
	assert.Nil(suite.T(), err, "Error creating HTTP request")
	response, err = client.Do(request)
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusNotFound, "Unexpected HTTP status code")

	response, err = http.Get(suite.makeRequestURL("/messages/2"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &getMessageResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.NotNil(suite.T(), getMessageResponse.Message, "Response did not include valid message")
	assert.Equal(suite.T(), getMessageResponse.Id, uint64(2), "Unexpected message id in response")
	assert.Equal(suite.T(), getMessageResponse.Payload, "foo", "Unexpected message payload response")

	message.Payload = "test"
	buf, err = json.Marshal(message)
	assert.Nil(suite.T(), err, "Error encoding json")
	for i := 0; i < 20; i++ {
		response, err = http.Post(suite.makeRequestURL("/messages"), "application/json", bytes.NewReader(buf))
		assert.Nil(suite.T(), err, "Error making HTTP request")
		assert.Equal(suite.T(), response.StatusCode, http.StatusCreated, "Unexpected HTTP status code")
	}

	response, err = http.Get(suite.makeRequestURL("/messages"))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	nextRelativeUrl := response.Header.Get("x-next-relative-url")
	assert.Equal(suite.T(), nextRelativeUrl, "/messages?afterId=21&limit=20", "Unexpected x-next-relative-url")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &listMessagesResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.Equal(suite.T(), len(listMessagesResponse), 20, "Unexpected number of messages in response")

	response, err = http.Get(suite.makeRequestURL(nextRelativeUrl))
	assert.Nil(suite.T(), err, "Error making HTTP request")
	assert.Equal(suite.T(), response.StatusCode, http.StatusOK, "Unexpected HTTP status code")
	assert.Equal(suite.T(), response.Header.Get("content-type"), "application/json", "Unexpected content-type")
	assert.Equal(suite.T(), response.Header.Get("x-next-relative-url"), "", "Unexpected x-next-relative-url")
	buf, err = ioutil.ReadAll(response.Body)
	assert.Nil(suite.T(), err, "Error reading response body")
	response.Body.Close()
	err = json.Unmarshal(buf, &listMessagesResponse)
	assert.Nil(suite.T(), err, "Error decoding json")
	assert.Equal(suite.T(), len(listMessagesResponse), 1, "Unexpected number of messages in response")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndTestSuite))
}
