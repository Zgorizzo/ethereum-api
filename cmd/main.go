package main

import (
	"fmt"
	"runtime"

	"github.com/INFURA/infra-test-benjamin-mateo/api"
	"github.com/INFURA/infra-test-benjamin-mateo/config"
	"github.com/INFURA/infra-test-benjamin-mateo/logger"
	"github.com/gorilla/mux"

	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// this is to avoid repeating exit code we return errors that are catched in main.
// Usefull when we have to handle lots of initialisation stuff (e.g. DB, connection to external services ..)
func run() error {
	// load configuration
	config.Load()

	// get an API server
	s := api.NewServer(logger.Init(config.ReadBool("ENABLE_DEBUG")), mux.NewRouter())

	servingURL := fmt.Sprintf("%s:%d", config.ReadString("APP_URL"), config.ReadInt("APP_PORT"))
	fmt.Printf("running on %v cpus\n", runtime.NumCPU())
	s.Serve(servingURL)

	return nil
}
