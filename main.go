package main

import (
    "errors"
    //"fmt"
    "log"
    "net/http"

    "github.com/brandonto/rest-api-microservice-demo/api"
    "github.com/brandonto/rest-api-microservice-demo/db"
)

const dbFile = "/home/brandonto/rest-api-microservice-demo.db"
const dbBucketName = "DetailedMessageBucket"
const serverAddr = ":12345"

func main() {
    // Create and configure Db
    //
    dbCfg := db.Config{dbFile, dbBucketName}
    svcDb := db.NewDb(dbCfg)

    // Just quit if Db initialization fails
    //
    if err := svcDb.Initialize(); err != nil {
        log.Fatal(errors.New("Unable to initialize DB"))
    }
    defer svcDb.Close()

    // Create the router and start accepting connections
    //
    router := api.NewRouter(svcDb)
    http.ListenAndServe(serverAddr, router)
}
