package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brandonto/rest-api-microservice-demo/core"
	"github.com/brandonto/rest-api-microservice-demo/db"
)

const (
	envDbBucketName = "BUCKET"
	envPortName     = "PORT"
	dbPathName      = "DB_PATH"
)

func getConfig() core.Config {
	envPort, cP := os.LookupEnv(envPortName)
	bucket, cB := os.LookupEnv(envDbBucketName)
	dbFile, cD := os.LookupEnv(dbPathName)
	if !cP || !cB || !cD {
		fmt.Errorf("enviroment variables: [%s, %s, %s] are required...", envDbBucketName, envPortName, dbPathName)
	}
	port, err := strconv.ParseUint(envPort, 10, 64)
	if err != nil {
		fmt.Errorf("PORT environment variable must be an unsigned integer, value was: %s", envPort)
	}

	if port < uint64(49152) || port > uint64(65535) {
		fmt.Errorf("PORT environment variable must be a value between 49152–65535, value was: %s", envPort)
	}

	coreCfg := core.Config{
		DbCfg: db.Config{
			FilePath:   dbFile,
			BucketName: bucket,
		},
		Port:         port,
		EnableLogger: true,
		Standalone:   true,
	}
	return coreCfg
}

func main() {
	core.Run(getConfig())
}
