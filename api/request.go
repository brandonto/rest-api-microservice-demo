package api

import (
    "github.com/brandonto/rest-api-microservice-demo/model"
)

// ListMessagesRequest
//
const ListMessagesLimitQueryParamDefault = uint64(20)
const ListMessagesLimitQueryParamMin = uint64(0)
const ListMessagesLimitQueryParamMax = uint64(100)

const ListMessagesAfterIdQueryParamDefault = uint64(0)

const ListMessagesDetailedQueryParamDefault = false


// CreateMessageRequest
//
type CreateMessageRequest struct {
    *model.Message
}


// GetMessageRequest
const GetMessageDetailedQueryParamDefault = false


// PutMessageRequest
//
type PutMessageRequest struct {
    *model.Message
}
