package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// config struct
type config struct {
	Port  int `yaml:"port"`
	Mongo struct {
		Host string `yaml:"host"`
		Name string `yaml:"name"`
		Mode string `yaml:"mode"`
	}
}

// Config is globle configuration
var Config *config

// Init read the configuration file
func Init(conf string) {
	data, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Panic(err.Error())
	}

	Config = &config{}

	err = yaml.Unmarshal([]byte(data), Config)
	if err != nil {
		log.Panic(err.Error())
	}
}
