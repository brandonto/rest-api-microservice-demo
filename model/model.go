package model

type Message struct {
    Id uint64 `json:"id"`
    Payload string `json:"payload"`
}

type MessageMetadata struct {
    Palindrome bool `json:"palindrome"`
}

type DetailedMessage struct {
    Message *Message `json:"message"`
    Metadata *MessageMetadata `json:"metadata"`
}

type ModelError struct {
    Code int32 `json:"code"`
    Message string `json:"message"`
}
