package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Mode  string `json:"mode"`
	Data  string `json:"data"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Elasticsearch struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Key  string `json:"key"`
	} `json:"elasticsearch"`
	Sqlite struct {
		Path string `json:"path"`
	} `json:"sqlite"`
	Prometheus struct {
		Route string `json:"route"`
		Port  int    `json:"port"`
	} `json:"prometheus"`
}

func LoadConfig() (*Config, error) {
	config := Config{}

	// determine config path to use
	path, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	// get config file ref
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// decode config json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("FATAL: LoadConfig: %s", err.Error())
	}
	return &config, nil
}

func getConfigPath() (string, error) {
	var path string
	configIndex := -1

	// walk arguments passed to program
	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--config" || arg == "-c" {
			configIndex = i
			break
		}
	}

	// --config | -c passed. parse value
	if configIndex >= 0 {
		if len(os.Args) > configIndex {
			path = os.Args[configIndex+1]
		} else {
			return "", errors.New("invalid config flag. usage: --config|-c \"<path-to=config>\"")
		}
	} else {
		// Set default config path
		path = "./_config.json"
	}

	// verify config exists
	_, err := os.Stat(path)
	if err != nil {
		fatal := fmt.Sprintf("config file not found. %s", path)
		return "", errors.New(fatal)
	}
	return path, nil
}
