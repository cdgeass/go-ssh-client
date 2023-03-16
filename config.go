package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	// server name
	Name string `yaml:"name"`

	// server address
	Host string `yaml:"host"`

	// server port
	Port int `yaml:"port"`

	// login username
	Username string `yaml:"username"`

	// login password
	Password string `yaml:"password"`

	// init command
	Cmd string `yaml:"cmd"`
}

type ServerGroup struct {
	// group name
	Name string

	// servers
	Servers []Server
}

func LoadConfig() Config {
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	config := Config{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatal("Failed to parse config: ", err)
	}
	return config
}
