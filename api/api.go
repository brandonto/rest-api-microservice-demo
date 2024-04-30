package api

import (
    //"fmt"
    "net/http"

    "github.com/brandonto/rest-api-microservice-demo/db"
    "github.com/brandonto/rest-api-microservice-demo/model"

    "github.com/go-chi/render"
)

func ListMessages(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("List"))
    }
}

func CreateMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        request := &CreateMessageRequest{}

        // Parse and validate the request
        //
        if err := render.Decode(r, request); err != nil {
            // Respond with status Bad Request - no response payload
            //
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        message := request.Message
        palindrome := isPalindrome(message.Payload)
        metadata := &model.MessageMetadata{Palindrome: palindrome}
        detailedMessage := &model.DetailedMessage{Message: message,
                                                  Metadata: metadata}
        svcDb.CreateMessage(detailedMessage)

        // Respond with status OK - no response payload
        //
        w.WriteHeader(http.StatusOK)
        return
    }
}

func GetMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Get"))
    }
}

func UpdateMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Update"))
    }
}

func DeleteMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Delete"))
    }
}
