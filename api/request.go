package api

import (
    "github.com/brandonto/rest-api-microservice-demo/model"
)

type CreateMessageRequest struct {
    *model.Message
}

type PutMessageRequest struct {
    *model.Message
}
