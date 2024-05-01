package api

import (
    //"fmt"
    "net/http"
    "strings"

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
        // ignored
        //
        message := request.Message
        palindrome := isPalindrome(message.Payload)
        metadata := &model.MessageMetadata{Palindrome: palindrome}
        detailedMessage := &model.DetailedMessage{Message: message, Metadata: metadata}

        // Adds the message to the database
        //
        if err := svcDb.CreateMessage(detailedMessage); err != nil {
            // Respond with status Unprocessable content - no response payload
            //
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        // Respond with status Created - no response payload
        //
        w.WriteHeader(http.StatusCreated)
        return
    }
}

func GetMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        detailed := false

        // A missing "detailed" query param is treated as it being false
        //
        if detailedQueryParam := r.URL.Query().Get("detailed"); detailedQueryParam != "" {
            detailedQueryParam = strings.ToLower(detailedQueryParam)
            if detailedQueryParam == "1" || detailedQueryParam == "true" {
                detailed = true
            } else if detailedQueryParam == "0" || detailedQueryParam == "false" {
                detailed = false
            } else {
                // Respond with status Bad Request - no response payload
                //
                // TODO Do we need a payload here to indicate that query param
                //      "detailed" is malformed?
                //
                w.WriteHeader(http.StatusBadRequest)
                return
            }
        }

        // Response with status OK - response payload depends on the "detailed"
        // query param
        //
        detailedMessage := r.Context().Value("detailedMessage").(*model.DetailedMessage)
        render.Status(r, http.StatusOK)
        if detailed {
            render.JSON(w, r, detailedMessage)
        } else {
            render.JSON(w, r, detailedMessage.Message)
        }
    }
}

func UpdateMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        request := &PutMessageRequest{}

        // Parse and validate the request
        //
        if err := render.Decode(r, request); err != nil {
            // Respond with status Bad Request - no response payload
            //
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // Replaces the payload of the Message stored in the database with the
        // payload of the Message received in the request. Also updates any
        // relevent metadata of the message while we're at it.
        //
        message := request.Message
        detailedMessage := r.Context().Value("detailedMessage").(*model.DetailedMessage)
        detailedMessage.Message.Payload = message.Payload
        detailedMessage.Metadata.Palindrome = isPalindrome(message.Payload)

        // Replaces the message in the database
        //
        if err := svcDb.UpdateMessage(detailedMessage); err != nil {
            // Respond with status Unprocessable content - no response payload
            //
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        // Respond with status No Content - no response payload
        //
        w.WriteHeader(http.StatusNoContent)
        return
    }
}

func DeleteMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        detailedMessage := r.Context().Value("detailedMessage").(*model.DetailedMessage)
        messageId := detailedMessage.Message.Id

        if err := svcDb.DeleteMessage(messageId); err != nil {
            // Respond with status Unprocessable content - no response payload
            //
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        // Respond with status OK - no response payload
        //
        w.WriteHeader(http.StatusCreated)
        return
    }
}
