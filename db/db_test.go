package db

import (
	"testing"

	"github.com/brandonto/rest-api-microservice-demo/model"

	"github.com/stretchr/testify/assert"
)

var testDbCfg = Config{
	FilePath:   "/home/brandonto/rest-api-microservice-demo.db",
	BucketName: "UnitTestBucket",
}

func TestBasicFunctionality(t *testing.T) {
	db := NewDb(testDbCfg)
	db.Initialize()
	defer db.Close()

	// Clear the database to ensure that we're in a clean state
	//
	assert.Nil(t, db.ClearMessages(), "ClearMessages() failed")

	// Try to retrieve a message that doesn't exist, this also checks if the
	// database is empty - shound fail
	//
	_, err := db.GetMessage(1)
	assert.NotNil(t, err, "There should not be a message returned from this call")

	// Create a message with id=6 (the id should be ignored and replaced with an
	// internal id))
	//
	message := &model.Message{Id: 6, Payload: "create"}
	metadata := &model.MessageMetadata{Palindrome: false}
	detailedMessage := &model.DetailedMessage{Message: message, Metadata: metadata}
	assert.Nil(t, db.CreateMessage(detailedMessage), "CreateMessage() failed")

	// Try to retrieve a message at id=6, should fail
	//
	_, err = db.GetMessage(6)
	assert.NotNil(t, err, "There should not be a message returned from this call")

	// Try to retrieve a message at id=1, should succeed
	//
	detailedMessageFromDb, err := db.GetMessage(1)
	assert.Nil(t, err, "There should be a message returned from this call")

	// Check if the message retrieved has it's id set to 1 in the payload
	//
	assert.Equal(t, uint64(1), detailedMessageFromDb.Message.Id, "Unexpected MessageId")

	// Check if the message retrieved has the correct payload
	//
	assert.Equal(t, "create", detailedMessageFromDb.Message.Payload, "Unexpected Payload")

	// Try to replace a message at id=1, should succeed
	//
	message.Id = 1
	message.Payload = "update"
	assert.Nil(t, db.UpdateMessage(detailedMessage), "UpdateMessage() failed")

	// Try to retrieve a message at id=1, should succeed
	//
	detailedMessageFromDb, err = db.GetMessage(1)
	assert.Nil(t, err, "There should be a message returned from this call")

	// Check if the message retrieved has the correct payload
	//
	assert.Equal(t, "update", detailedMessageFromDb.Message.Payload, "Unexpected Payload")

	// Try to delete a message at id=1, should succeed
	//
	assert.Nil(t, db.DeleteMessage(1), "DeleteMessage() failed")

	// Try to retrieve a message at id=1, should fail
	//
	detailedMessageFromDb, err = db.GetMessage(1)
	assert.NotNil(t, err, "There should not be a message returned from this call")

	// Try to create a message at id=1, should succeed but id should be ignored
	//
	message.Id = 1
	message.Payload = "create"
	assert.Nil(t, db.CreateMessage(detailedMessage), "CreateMessage() failed")

	// Try to retrieve a message at id=1, should fail (the id from the create was
	// ignored)
	//
	detailedMessageFromDb, err = db.GetMessage(1)
	assert.NotNil(t, err, "There should not be a message returned from this call")

	// Try to retrieve a message at id=2, should succeed
	//
	detailedMessageFromDb, err = db.GetMessage(2)
	assert.Nil(t, err, "There should be a message returned from this call")

	// Check if the message retrieved has the correct payload
	//
	assert.Equal(t, "create", detailedMessageFromDb.Message.Payload, "Unexpected Payload")

	// Create 40 additional messages, should all succeed but id should be ignored
	//
	for i := 0; i < 41; i++ {
		assert.Nil(t, db.CreateMessage(detailedMessage), "CreateMessage() failed")
	}

	// Get a list of 20 messages starting AFTER message id=2
	//
	detailedList, afterId, err := db.ListMessages(20, 2+1)
	assert.Nil(t, err, "ListMessages() failed")
	assert.Equal(t, afterId, uint64(22), "unexpected afterId")
	assert.Equal(t, len(detailedList), 20, "unexpected list size")

	// Get a list of the next 20 messages starting after message id=22
	//
	detailedList, afterId, err = db.ListMessages(20, afterId+1)
	assert.Nil(t, err, "ListMessages() failed")
	assert.Equal(t, afterId, uint64(42), "unexpected afterId")
	assert.Equal(t, len(detailedList), 20, "unexpected list size")

	// Get a list of the next 20 messages starting after message id=42.
	// Since there is only 1 element left on the list prior to this call, the
	// afterId returned should be 0 to indicate that we're finished iterating
	// through the database.
	//
	detailedList, afterId, err = db.ListMessages(20, afterId+1)
	assert.Nil(t, err, "ListMessages() failed")
	assert.Equal(t, afterId, uint64(0), "unexpected afterId")
	assert.Equal(t, len(detailedList), 1, "unexpected list size")
}
