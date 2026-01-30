package main

import (
	"flag"
	"time"

	"github.com/pisue/go-playground/grpc-server/cmd"
	"github.com/pisue/go-playground/grpc-server/config"
	"github.com/pisue/go-playground/grpc-server/gRPC/server"
)

var configFlag = flag.String("config", "./config.toml", "config path")

func main() {
	flag.Parse()
	cfg := config.NewConfig(*configFlag)

	err := server.NewGRPCServer(cfg)
	if err != nil {
		panic(err)
	}

	time.Sleep(1e9)
	cmd.NewApp(cfg)
}
