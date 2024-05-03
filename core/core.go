package core

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/brandonto/rest-api-microservice-demo/api"
	"github.com/brandonto/rest-api-microservice-demo/db"
)

type Config struct {
	DbCfg        db.Config
	Port         uint64
	EnableLogger bool
	Standalone   bool
}

const ServerShutdownTimeoutInSeconds = 3

func Run(coreCfg Config) {
	var err error

	// Create and configure Db
	//
	svcDb := db.NewDb(coreCfg.DbCfg)

	// Just quit if Db initialization fails
	//
	if err = svcDb.Initialize(); err != nil {
		log.Fatal(errors.New("Unable to initialize DB"))
	}
	defer svcDb.Close()

	// Set up HTTP routes
	//
	router := api.NewRouter(svcDb, coreCfg.EnableLogger, coreCfg.Standalone)

	// Create and configure the server and start accepting connections
	//
	portStr := strconv.FormatUint(coreCfg.Port, 10)
	server := &http.Server{Addr: "localhost:" + portStr, Handler: router}

	// Start server in goroutine so we can attempt to perform a graceful shutdown
	// when sent a signal. Most commonly SIGTERM and SIGINT when running the
	// server as a standalone app. But perhaps more importantly, SIGUSR1 when
	// running the server as part of the e2e test suite for set up and tear down
	// needs in between tests... for example, if we need to shut down the server
	// in order to clear the database in between tests.
	//
	// https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
	// https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	//
	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGUSR1)
	if coreCfg.Standalone {
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	}
	<-sigChan

	// Small timeout to prevent http.ShutDown from hanging indefinitely if any
	// client can't be gracefully shut down.
	//
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), ServerShutdownTimeoutInSeconds*time.Second)
	defer shutdownRelease()

	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
