package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config contains all settings
type Config struct {
	Database DBConfig   `yaml:"db"`
	View     ViewConfig `yaml:"view"`
}

// ReadFile loads a configuration from file
func ReadFile(filename string) *Config {
	log.Printf("Reading configuration file '%s'", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	conf := new(Config)
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		log.Fatalln(err)
	}

	return conf
}
