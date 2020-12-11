package main

import (
	"simian/server"
	"simian/service"

	sample "simian/sample"
)

func main() {
	server.SetCORS()

	server.AddServices(service.Collection{
		sample.Config,
	})

	server.Start()
}
