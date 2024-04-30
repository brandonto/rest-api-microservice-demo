package main

import (
    //"fmt"
    "net/http"

    "github.com/brandonto/rest-api-microservice-demo/api"
    "github.com/brandonto/rest-api-microservice-demo/db"
)

func main() {
    svcDb := db.Open("/home/brandonto/rest-api-microservice-demo.db")
    defer svcDb.Close()

    router := api.CreateRouter(svcDb)
    http.ListenAndServe(":12345", router)
}
