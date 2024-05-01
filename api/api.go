package api

import (
    //"fmt"
    "net/http"
    "strconv"

    "github.com/brandonto/rest-api-microservice-demo/db"
    "github.com/brandonto/rest-api-microservice-demo/model"

    "github.com/go-chi/render"
)

func ListMessages(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        // Retrieve pagination query params from context
        //
        detailed := r.Context().Value("detailed").(bool)
        limit := r.Context().Value("limit").(uint64)
        afterId := r.Context().Value("afterId").(uint64)

        detailedMessages, nextAfterId, err := svcDb.ListMessages(limit, afterId + 1)
        if err != nil {
            // Something went wrong with a batch get... respond with status
            // Unprocessable content - no response payload
            //
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        // If there are further message to retrieve from the database, we'll
        // return a relative URL for the next page of messages in an HTTP header
        //
        if nextAfterId != 0 {
            // Start with the afterId
            //
            nextRelativeUrl := r.URL.Path + "?afterId=" + strconv.FormatUint(nextAfterId, 10)

            // "limit" will always have a value, we assign a default value if
            // one wasn't included in the HTTP request
            //
            nextRelativeUrl = nextRelativeUrl + "+limit=" + strconv.FormatUint(limit, 10)

            // "detailed" defaults to false, so we add the query param only if
            // it was set to true
            //
            if detailed {
                nextRelativeUrl += "+detailed=true"
            }
	        w.Header().Set("x-next-relative-url", nextRelativeUrl)
        }

        // Response with status OK - response payload depends on the "detailed"
        // query param
        //
        render.Status(r, http.StatusOK)
        if detailed {
            render.JSON(w, r, detailedMessages)
        } else {
            // Transform array of DetailedMessages into array of Messages
            //
            var messages []*model.Message
            for _, v := range detailedMessages {
                messages = append(messages, v.Message)
            }
            render.JSON(w, r, messages)
        }
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
        // ignored in the database layer
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
        var err error

        // A missing "detailed" query param is treated as it being false
        //
        detailed := GetMessageDetailedQueryParamDefault
        if detailedQueryParam := r.URL.Query().Get("detailed"); detailedQueryParam != "" {
            detailed, err = stringToBool(detailedQueryParam)
            if err != nil {
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

        // Respond with status No Content - no response payload
        //
        w.WriteHeader(http.StatusNoContent)
        return
    }
}
