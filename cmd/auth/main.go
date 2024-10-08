package main

import (
	"github.com/emaforlin/auth-service/internal/config"
	"github.com/emaforlin/auth-service/internal/server"
	"github.com/hashicorp/go-hclog"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Debug,
	})

	config.InitViper("config")
	conf := config.LoadConfig()
	server.NewRPCServer(logger, conf).Start()
}
