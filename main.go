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
