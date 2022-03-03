package internal

import (
	"encoding/json"
	"os"
)

type Config struct {
	ZoneName string `json:"ZoneName"`
	ZoneId   string `json:"ZoneId"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func SplitLongRoute53String(input string) string {
	i := 250
	output := input
	for i <= len(output) {
		output = output[:i] + "\"\"" + output[i:]
		i += 252
	}
	return output
}
