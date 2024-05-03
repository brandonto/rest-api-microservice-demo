package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/brandonto/rest-api-microservice-demo/core"
	"github.com/brandonto/rest-api-microservice-demo/db"
)

const defaultDbBucketName = "DetailedMessageBucket"
const defaultPort = uint64(55555)

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
	port := defaultPort
	if len(os.Args) > 2 {
		// Some sanity check for the port argument before using
		//
		var err error
		port, err = strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("port argument must be an unsigned integer")
			return
		}

		if port < uint64(49152) || port > uint64(65535) {
			fmt.Println("port argument must be a value between 49152–65535")
			return
		}
	}

	// Third (optional) argument is the the bucket name
	//
	dbBucketName := defaultDbBucketName
	if len(os.Args) > 3 {
		dbBucketName = os.Args[3]
	}

	// Configure and run the application
	//
	dbCfg := db.Config{
		FilePath:   dbFile,
		BucketName: dbBucketName,
	}

	coreCfg := core.Config{
		DbCfg:        dbCfg,
		Port:         port,
		EnableLogger: true,
		Standalone:   true,
	}

	core.Run(coreCfg)
}
