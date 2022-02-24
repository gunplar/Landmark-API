package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	ZoneName string `json:"ZoneName"`
	ZoneId   string `json:"ZoneId"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer func(configFile *os.File) {
		err := configFile.Close()
		check(err)
	}(configFile)
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	check(err)
	return config
}
