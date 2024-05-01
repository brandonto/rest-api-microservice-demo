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

        // Transform the Message received in the request into a DetailedMessage
        // to store in the database. The "id" field of the message will be
        // ignoredk
        //
        message := request.Message
        palindrome := isPalindrome(message.Payload)
        metadata := &model.MessageMetadata{Palindrome: palindrome}
        detailedMessage := &model.DetailedMessage{Message: message,
                                                  Metadata: metadata}

        // Adds the message to the database
        //
        if err := svcDb.CreateMessage(detailedMessage); err != nil {
            // Respond with status Unprocessable content - no response payload
            //
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        // Respond with status OK - no response payload
        //
        w.WriteHeader(http.StatusOK)
        return
    }
}

func GetMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        //w.Write([]byte("Get"))
        message := r.Context().Value("message").(*model.Message)
        metadata := &model.MessageMetadata{Palindrome: true}
        detailedMessage := model.DetailedMessage{Message: message, Metadata: metadata}

        //render.JSON(w, r, message)
        render.JSON(w, r, detailedMessage)
        //if err := render.JSON(w, r, message); err != nil {
        //    w.WriteHeader(http.StatusUnprocessableEntity)
        //    return
        //}
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
