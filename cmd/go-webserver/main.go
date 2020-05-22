package main

import (
	"os"

	"github.com/pieterclaerhout/go-log"
	webserver "github.com/pieterclaerhout/go-webserver/v2"
)

func main() {

	// Setup logging
	log.PrintColors = true
	log.PrintTimestamp = true
	log.DebugMode = (os.Getenv("DEBUG") == "1")

	// Run the app with the server
	err := webserver.New().RunWithApps(
		&SampleApp{},
	)
	log.CheckError(err)

}
