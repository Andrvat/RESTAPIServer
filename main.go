package main

import (
	"awesomeProject/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configDefaultPath string
)

func init() {
	flag.StringVar(&configDefaultPath,
		"config-path",
		"configs/apiserver.toml",
		"Initialize path to config TOML file")
}

func main() {
	flag.Parse()
	config := apiserver.NewDefaultConfig()
	if _, err := toml.DecodeFile(configDefaultPath, config); err != nil {
		log.Fatal(err)
	}
	server := apiserver.NewServer(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
