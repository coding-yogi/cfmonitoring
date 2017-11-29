package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Settings ...
type Settings struct {
	RabbitMQ RabbitMQ
	InfluxDB InfluxDB
	Cf       Cf
}

// RabbitMQ ...
type RabbitMQ struct {
	Host   string
	Port   int32
	Queues []string
}

//InfluxDB ...
type InfluxDB struct {
	Host     string
	Port     int32
	Database string
}

//CF
type Cf struct {
	API   string
	Org   string
	Space string
	Apps  []string
}

// GetConfig ...
func GetConfig() (Settings, error) {
	var settings Settings

	pwd, _ := os.Getwd()
	file, err := ioutil.ReadFile(pwd + "/config.json")
	if err != nil {
		return settings, err
	}

	settings = Settings{}
	json.Unmarshal(file, &settings)

	return settings, nil
}
