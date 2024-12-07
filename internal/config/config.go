package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Mode string `json:"mode"`
	Api  struct {
		Address     string `json:"address"`
		WriteTimout int    `json:"write_timeout"`
		ReadTimout  int    `json:"read_timeout"`
		IdleTimeout int    `json:"idle_timeout"`
	} `json:"api"`
	Data  string `json:"data"`
	Redis struct {
		Host                string `json:"host"`
		Port                int    `json:"port"`
		DefaultCacheSeconds int    `json:"default_cache_seconds"`
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
		Route   string `json:"route"`
		Address string `json:"address"`
	} `json:"prometheus"`
}

func NewConfig() (*Config, error) {
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
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}
	return &config, nil
}

func getConfigPath() (string, error) {
	var path string

	defaultConfigPaths := []string{"./", "./configs/"}
	defaultConfigFile := "_config.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	// no override config
	if path == "" {
		// Check for default config file in running dir
		for _, defaultConfigPath := range defaultConfigPaths {
			filePath := defaultConfigPath + defaultConfigFile
			_, err := os.Stat(filePath)
			if err == nil {
				path = filePath
				break
			}
		}
	}

	// verify config exists
	_, err := os.Stat(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("config file not found. %s", path))
	}
	return path, nil
}
