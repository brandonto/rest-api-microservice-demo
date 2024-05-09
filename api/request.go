package api

import (
	"errors"
	"net/http"

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

func (decodedReq *CreateMessageRequest) Bind(r *http.Request) error {
	if decodedReq.Message == nil || decodedReq.Message.Payload == "" {
		return errors.New("Missing required fields")
	}

	return nil
}

// GetMessageRequest
//
const GetMessageDetailedQueryParamDefault = false

// PutMessageRequest
//
type PutMessageRequest struct {
	*model.Message
}

func (decodedReq *PutMessageRequest) Bind(r *http.Request) error {
	if decodedReq.Message == nil || decodedReq.Message.Payload == "" {
		return errors.New("Missing required fields")
	}

	return nil
}
