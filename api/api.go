package api

import (
    "net/http"

    "github.com/brandonto/rest-api-microservice-demo/db"
    //"github.com/brandonto/rest-api-microservice-demo/model"

    //"github.com/go-chi/render"
)

func ListMessages(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("List"))
    }
}

func CreateMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Create"))
    }
}

func GetMessage(svcDb *db.Db) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Get"))
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
