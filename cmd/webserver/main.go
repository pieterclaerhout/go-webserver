package main

import (
	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver"
	"github.com/pieterclaerhout/go-webserver/cmd/webserver/core"
)

func main() {

	server := webserver.New()

	server.Register(&core.Core{})

	err := server.Start()
	log.CheckError(err)

}
