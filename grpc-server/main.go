package main

import (
	"flag"

	"github.com/pisue/go-playground/grpc-server/cmd"
	"github.com/pisue/go-playground/grpc-server/config"
)

var configFlag = flag.String("config", "./config.toml", "config path")

func main() {
	flag.Parse()
	cfg := config.NewConfig(*configFlag)
	cmd.NewApp(cfg)
}
