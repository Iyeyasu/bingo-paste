package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Configuration contains all settings
type Configuration struct {
	Database DBConfiguration   `yaml:"db"`
	View     ViewConfiguration `yaml:"view"`
}

// ReadFile loads a configuration from file
func ReadFile(filename string) *Configuration {
	log.Printf("Reading configuration file '%s'", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	conf := new(Configuration)
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		log.Fatalln(err)
	}

	return conf
}
