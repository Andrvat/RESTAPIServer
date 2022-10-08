package main

import (
	"awesomeProject/internal/app/apiserver"
	"flag"
	"log"
)

var (
	serverConfigPath string
)

func init() {
	flag.StringVar(&serverConfigPath,
		"config-path",
		"configs/apiserver.toml",
		"Initialize path to config TOML file")
}

// @title Test API server
// @version 0.1
// @description API Server for learning Go lang

// @host localhost::5544
// @BasePath /

func main() {
	flag.Parse()
	config, err := apiserver.NewConfigFromToml(serverConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
