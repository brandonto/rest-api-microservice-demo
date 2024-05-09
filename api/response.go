package api

import (
	"github.com/brandonto/rest-api-microservice-demo/model"
)

// ListMessagesResponse
//
type ListMessagesResponse []model.Message

// GetMessageResponse
//
type GetMessageResponse struct {
	*model.Message
}
