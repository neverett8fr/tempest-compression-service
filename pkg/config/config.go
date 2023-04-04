package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	route = "config/config.yaml"
)

type Service struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DB struct {
	Driver   string `yaml:"driver"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type Config struct {
	Service Service `yaml:"user-service-config"`
	DB      DB      `yaml:"user-db-config"`
}

func Initialise() (*Config, error) {

	conf := Config{}

	yamlFile, err := ioutil.ReadFile(route)
	if err != nil {
		return &Config{}, fmt.Errorf("issue finding config yaml, err %v", err)
	}
	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return &Config{}, fmt.Errorf("issue unmarshalling config yaml, err %v", err)
	}

	return &conf, nil
}
