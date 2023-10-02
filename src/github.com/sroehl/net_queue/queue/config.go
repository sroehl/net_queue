package queue

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port       int    `yaml:port`
		ListenHost string `yaml:listenhost`
	} `yaml:server`
}

func Read_config(filename string) (Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func New_config(host string, port int) Config {
	cfg := Config{}
	cfg.Server.ListenHost = host
	cfg.Server.Port = port
	return cfg
}
