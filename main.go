package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/brandonto/rest-api-microservice-demo/api"
    "github.com/brandonto/rest-api-microservice-demo/db"
)

const defaultDbBucketName = "DetailedMessageBucket"
const defaultPortStr = "55555"

func main() {
    // Outputs usage if number of arguments are off
    //
    if len(os.Args) > 4 || len(os.Args) < 2 {
        fmt.Println("usage: " + os.Args[0] + " db_path [port] [db_bucket_name]")
        return
    }

    // First argument is the filepath of the database
    //
    dbFile := os.Args[1]

    // Second (optional) argument is the the port number
    //
    portStr := defaultPortStr
    if len(os.Args) > 2 {
        // Some sanity check for the port argument before using
        //
        port, err := strconv.ParseUint(os.Args[2], 10, 64)
        if err != nil {
            fmt.Println("port argument must be an unsigned integer")
            return
        }

        if port < uint64(49152) || port > uint64(65535) {
            fmt.Println("port argument must be a value between 49152–65535")
            return
        }

        portStr = os.Args[2]
    }

    // Third (optional) argument is the the bucket name
    //
    dbBucketName := defaultDbBucketName
    if len(os.Args) > 3 {
        dbBucketName = os.Args[3]
    }

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
    http.ListenAndServe(":" + portStr, router)
}
