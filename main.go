package main

import (
	"awesomeProject/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
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
	config := apiserver.NewDefaultConfig()
	if _, err := toml.DecodeFile(serverConfigPath, config); err != nil {
		log.Fatal(err)
	}
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
