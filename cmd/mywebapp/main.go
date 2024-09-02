package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/VeeRomanoff/mywebapp/internal/mywebapp"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "conf", "configs/api.toml", "configuration file path")
}

func main() {
	flag.Parse()
	log.Println("app started")

	// setting up configuration
	config := mywebapp.NewConfig()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		log.Fatal("failed to load configuration. using default configuration")
	}
	log.Println("config database_api", config.StorageConfig.DatabaseURI)

	// setting up server
	server := mywebapp.NewMyWebApp(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
