package core

import (
    "errors"
    "log"
    "net/http"
    "strconv"

    "github.com/brandonto/rest-api-microservice-demo/api"
    "github.com/brandonto/rest-api-microservice-demo/db"
)

type Config struct {
    DbCfg db.Config
    Port uint64
    EnableLogger bool
}

func Run(coreCfg Config) {
    // Create and configure Db
    //
    svcDb := db.NewDb(coreCfg.DbCfg)

    // Just quit if Db initialization fails
    //
    if err := svcDb.Initialize(); err != nil {
        log.Fatal(errors.New("Unable to initialize DB"))
    }
    defer svcDb.Close()

    // Create the router and start accepting connections
    //
    router := api.NewRouter(svcDb, coreCfg.EnableLogger)
    portStr := strconv.FormatUint(coreCfg.Port, 10)
    http.ListenAndServe(":" + portStr, router)
}
