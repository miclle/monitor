package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// MongoConfig mongodb configuration struct
type MongoConfig struct {
	Host string `yaml:"host"`
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

// config struct
type config struct {
	Port  int         `yaml:"port"`
	Mongo MongoConfig `yaml:"mongo"`
}

// Port is app server port
var Port int

// Mongo configuration
var Mongo MongoConfig

// init read the configuration file
func init() {
	config := &config{
		Port: 8080,
		Mongo: MongoConfig{
			Host: "127.0.0.1:27017",
			Name: "observer",
			Mode: "strong",
		},
	}

	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Print("Connot found the config.yml file, use default configuration!")
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Panic(err.Error())
	}

	Port = config.Port
	Mongo = config.Mongo
}
