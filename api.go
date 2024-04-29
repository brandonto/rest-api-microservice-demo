package main

import (
    "net/http"
)

func ListMessages(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("List"))
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Create"))
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Get"))
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Update"))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Delete"))
}
