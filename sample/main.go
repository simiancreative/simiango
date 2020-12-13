package main

import (
	_ "github.com/simiancreative/simiango/config"

	"github.com/simiancreative/simiango/data/redis"
	"github.com/simiancreative/simiango/server"
	"github.com/simiancreative/simiango/service"

	sample "github.com/simiancreative/simiango/sample/sample"
)

func main() {
	redis.Connect()
	server.SetCORS()

	server.AddServices(service.Collection{
		sample.Config,
	})

	server.Start()
}
