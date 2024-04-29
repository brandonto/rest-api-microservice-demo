package main

type Message struct {
    Id int64 `json:"id"`
    Payload string `json:"payload"`
}

type DetailedMessageMetadata struct {
    Palindrome bool `json:"palindrome"`
}

type DetailedMessage struct {
    Message *Message `json:"message"`
    Metadata *DetailedMessageMetadata `json:"metadata"`
}

type ModelError struct {
    Code int32 `json:"code"`
    Message string `json:"message"`
}
